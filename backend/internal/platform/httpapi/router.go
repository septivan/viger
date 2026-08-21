package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	gameservices "github.com/septivan/viger/backend/internal/core/game/services"
	review "github.com/septivan/viger/backend/internal/core/review/domain"
	reviewservices "github.com/septivan/viger/backend/internal/core/review/services"
	"github.com/septivan/viger/backend/internal/platform/observability"
)

const maximumReviewBody = 16 << 10

// Settings contains the explicit dependencies required by the HTTP adapter.
type Settings struct {
	Games            gameservices.Service
	Reviews          reviewservices.Service
	WebSocket        http.Handler
	Metrics          *observability.Metrics
	Logger           *slog.Logger
	AllowedOrigins   []string
	ReviewRateLimit  int
	ReviewRateWindow time.Duration
}

type api struct {
	settings Settings
	limiter  *rateLimiter
}

// NewRouter maps REST endpoints to application services and shared middleware.
func NewRouter(settings Settings) http.Handler {
	server := &api{settings: settings, limiter: newRateLimiter(settings.ReviewRateLimit, settings.ReviewRateWindow)}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", server.live)
	mux.HandleFunc("GET /health/ready", server.ready)
	mux.Handle("GET /metrics", settings.Metrics)
	mux.HandleFunc("GET /v1/games", server.listGames)
	mux.HandleFunc("GET /v1/games/{gameID}", server.getGame)
	mux.HandleFunc("GET /v1/games/{gameID}/reviews", server.listReviews)
	mux.HandleFunc("POST /v1/games/{gameID}/reviews", server.createReview)
	mux.Handle("GET /v1/ws", settings.WebSocket)
	return server.recoverPanic(server.requestLog(server.securityHeaders(server.cors(mux))))
}

func (server *api) live(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (server *api) ready(writer http.ResponseWriter, request *http.Request) {
	if _, err := server.settings.Games.List(request.Context(), gameservices.ListQuery{Page: 1, PageSize: 1}); err != nil {
		writeError(writer, request, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "The game catalog is unavailable.", nil)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (server *api) listGames(writer http.ResponseWriter, request *http.Request) {
	query, err := parseGameQuery(request)
	if err != nil {
		writeError(writer, request, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	page, err := server.settings.Games.List(request.Context(), query)
	if err != nil {
		server.writeServiceError(writer, request, err)
		return
	}
	items := make([]gameResponse, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, gameDTO(item))
	}
	writeJSON(writer, http.StatusOK, pageResponse[gameResponse]{Data: items, Pagination: paginationDTO(page.Page, page.PageSize, page.TotalItems, page.TotalPages)})
}

func (server *api) getGame(writer http.ResponseWriter, request *http.Request) {
	item, err := server.settings.Games.Find(request.Context(), request.PathValue("gameID"))
	if err != nil {
		server.writeServiceError(writer, request, err)
		return
	}
	writeJSON(writer, http.StatusOK, dataResponse[gameResponse]{Data: gameDTO(item)})
}

func (server *api) listReviews(writer http.ResponseWriter, request *http.Request) {
	page, pageSize, err := parsePage(request, 10)
	if err != nil {
		writeError(writer, request, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	result, err := server.settings.Games.ListReviews(request.Context(), request.PathValue("gameID"), gameservices.ReviewQuery{Sort: request.URL.Query().Get("sort"), Page: page, PageSize: pageSize})
	if err != nil {
		server.writeServiceError(writer, request, err)
		return
	}
	items := make([]reviewResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, reviewDTO(item))
	}
	writeJSON(writer, http.StatusOK, pageResponse[reviewResponse]{Data: items, Pagination: paginationDTO(result.Page, result.PageSize, result.TotalItems, result.TotalPages)})
}

func (server *api) createReview(writer http.ResponseWriter, request *http.Request) {
	if !server.limiter.allow(request) {
		writer.Header().Set("Retry-After", strconv.Itoa(int(server.settings.ReviewRateWindow.Seconds())))
		writeError(writer, request, http.StatusTooManyRequests, "RATE_LIMITED", "Too many reviews were submitted. Please try again later.", nil)
		return
	}
	if mediaType := request.Header.Get("Content-Type"); !strings.HasPrefix(strings.ToLower(mediaType), "application/json") {
		writeError(writer, request, http.StatusUnsupportedMediaType, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json.", nil)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, maximumReviewBody)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	var input struct {
		ReviewerName string `json:"reviewerName"`
		Rating       int    `json:"rating"`
		Text         string `json:"text"`
	}
	if err := decoder.Decode(&input); err != nil {
		message := "The request body must be a valid JSON object with known fields."
		if errors.As(err, new(*http.MaxBytesError)) {
			message = "The request body exceeds the 16 KiB limit."
		}
		writeError(writer, request, http.StatusBadRequest, "INVALID_JSON", message, nil)
		return
	}
	if err := ensureJSONEnd(decoder); err != nil {
		writeError(writer, request, http.StatusBadRequest, "INVALID_JSON", "The request body must contain exactly one JSON object.", nil)
		return
	}
	created, err := server.settings.Reviews.Create(request.Context(), reviewservices.CreateInput{
		GameID: request.PathValue("gameID"), ReviewerName: input.ReviewerName, Rating: input.Rating, Text: input.Text,
	})
	if err != nil {
		server.writeServiceError(writer, request, err)
		return
	}
	server.settings.Metrics.RecordReviewCreated()
	writer.Header().Set("Location", "/v1/games/"+created.GameID+"/reviews/"+created.ID)
	writeJSON(writer, http.StatusCreated, dataResponse[reviewResponse]{Data: reviewDTO(created)})
}

func parseGameQuery(request *http.Request) (gameservices.ListQuery, error) {
	page, pageSize, err := parsePage(request, 12)
	if err != nil {
		return gameservices.ListQuery{}, err
	}
	minimumRating := 0.0
	if raw := request.URL.Query().Get("minRating"); raw != "" {
		minimumRating, err = strconv.ParseFloat(raw, 64)
		if err != nil {
			return gameservices.ListQuery{}, fmt.Errorf("minRating must be a number between 0 and 5")
		}
	}
	return gameservices.ListQuery{
		Search: request.URL.Query().Get("q"), Genre: request.URL.Query().Get("genre"),
		Platform: request.URL.Query().Get("platform"), MinRating: minimumRating,
		Sort: request.URL.Query().Get("sort"), Page: page, PageSize: pageSize,
	}, nil
}

func parsePage(request *http.Request, defaultSize int) (int, int, error) {
	page, pageSize := 1, defaultSize
	var err error
	if raw := request.URL.Query().Get("page"); raw != "" {
		page, err = strconv.Atoi(raw)
		if err != nil {
			return 0, 0, fmt.Errorf("page must be a positive integer")
		}
	}
	if raw := request.URL.Query().Get("pageSize"); raw != "" {
		pageSize, err = strconv.Atoi(raw)
		if err != nil {
			return 0, 0, fmt.Errorf("pageSize must be an integer between 1 and 50")
		}
	}
	return page, pageSize, nil
}

func ensureJSONEnd(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("additional JSON value")
	}
	return nil
}

func (server *api) writeServiceError(writer http.ResponseWriter, request *http.Request, err error) {
	switch {
	case errors.Is(err, gameservices.ErrGameNotFound), errors.Is(err, reviewservices.ErrGameNotFound):
		writeError(writer, request, http.StatusNotFound, "GAME_NOT_FOUND", "The requested game was not found.", nil)
	case errors.Is(err, gameservices.ErrInvalidPage), errors.Is(err, gameservices.ErrInvalidPageSize), errors.Is(err, gameservices.ErrInvalidSearch), errors.Is(err, gameservices.ErrInvalidMinRating), errors.Is(err, gameservices.ErrInvalidGameSort), errors.Is(err, gameservices.ErrInvalidReviewSort):
		writeError(writer, request, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
	case errors.Is(err, review.ErrInvalidReviewerName):
		writeError(writer, request, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "The review contains invalid fields.", map[string]string{"reviewerName": err.Error()})
	case errors.Is(err, review.ErrInvalidRating):
		writeError(writer, request, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "The review contains invalid fields.", map[string]string{"rating": err.Error()})
	case errors.Is(err, review.ErrInvalidText), errors.Is(err, review.ErrControlCharacters):
		writeError(writer, request, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "The review contains invalid fields.", map[string]string{"text": err.Error()})
	default:
		server.settings.Logger.Error("request failed", "request_id", requestID(request.Context()), "error", err)
		writeError(writer, request, http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed.", nil)
	}
}

type dataResponse[T any] struct {
	Data T `json:"data"`
}

type pageResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination pagination `json:"pagination"`
}

type pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

type gameResponse struct {
	ID                 string         `json:"id"`
	Title              string         `json:"title"`
	Description        string         `json:"description"`
	Genre              string         `json:"genre"`
	Platforms          []string       `json:"platforms"`
	Developer          string         `json:"developer"`
	ReleaseDate        string         `json:"releaseDate"`
	AverageRating      float64        `json:"averageRating"`
	ReviewCount        int            `json:"reviewCount"`
	RatingDistribution map[string]int `json:"ratingDistribution"`
}

type reviewResponse struct {
	ID           string    `json:"id"`
	GameID       string    `json:"gameId"`
	ReviewerName string    `json:"reviewerName"`
	Rating       int       `json:"rating"`
	Text         string    `json:"text"`
	CreatedAt    time.Time `json:"createdAt"`
}

type errorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"requestId"`
}

func gameDTO(item gameservices.CatalogItem) gameResponse {
	distribution := make(map[string]int, 5)
	for rating := 1; rating <= 5; rating++ {
		distribution[strconv.Itoa(rating)] = item.Ratings.Distribution[rating]
	}
	return gameResponse{
		ID: item.Game.ID, Title: item.Game.Title, Description: item.Game.Description,
		Genre: item.Game.Genre, Platforms: append([]string(nil), item.Game.Platforms...),
		Developer: item.Game.Developer, ReleaseDate: item.Game.ReleaseDate.Format(time.DateOnly),
		AverageRating: item.Ratings.Average, ReviewCount: item.Ratings.Total, RatingDistribution: distribution,
	}
}

func reviewDTO(item review.Review) reviewResponse {
	return reviewResponse{ID: item.ID, GameID: item.GameID, ReviewerName: item.ReviewerName, Rating: item.Rating, Text: item.Text, CreatedAt: item.CreatedAt}
}

func paginationDTO(page, pageSize, totalItems, totalPages int) pagination {
	return pagination{Page: page, PageSize: pageSize, TotalItems: totalItems, TotalPages: totalPages}
}

func writeJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("Cache-Control", "no-store")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}

func writeError(writer http.ResponseWriter, request *http.Request, status int, code, message string, fields map[string]string) {
	writeJSON(writer, status, errorResponse{Error: apiError{Code: code, Message: message, Fields: fields, RequestID: requestID(request.Context())}})
}
