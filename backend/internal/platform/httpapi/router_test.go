package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/septivan/viger/backend/internal/adapters/outbound/memory"
	gameservices "github.com/septivan/viger/backend/internal/core/game/services"
	reviewservices "github.com/septivan/viger/backend/internal/core/review/services"
	"github.com/septivan/viger/backend/internal/platform/observability"
	"github.com/septivan/viger/backend/internal/platform/realtime"
)

func testRouter(t *testing.T, origin string, rateLimit int) (http.Handler, *realtime.Hub) {
	t.Helper()
	games, reviews, err := memory.Seed()
	if err != nil {
		t.Fatal(err)
	}
	store, err := memory.New(games, reviews)
	if err != nil {
		t.Fatal(err)
	}
	metrics := &observability.Metrics{}
	hub := realtime.NewHub([]string{origin}, 10, metrics)
	return NewRouter(Settings{
		Games:     gameservices.New(store, store),
		Reviews:   reviewservices.New(store, store, hub, reviewservices.SystemClock{}, reviewservices.RandomIDGenerator{}),
		WebSocket: hub, Metrics: metrics, Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		AllowedOrigins: []string{origin}, ReviewRateLimit: rateLimit, ReviewRateWindow: time.Minute,
	}), hub
}

func TestGameEndpointsReturnPaginatedCatalogAndDetail(t *testing.T) {
	router, _ := testRouter(t, "http://localhost:3000", 10)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/v1/games?q=hades&page=1&pageSize=12", nil))
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"title":"Hades"`) || !strings.Contains(response.Body.String(), `"totalItems":1`) {
		t.Fatalf("unexpected list response: status=%d body=%s", response.Code, response.Body.String())
	}
	response = httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/v1/games/game-002", nil))
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"ratingDistribution"`) {
		t.Fatalf("unexpected detail response: status=%d body=%s", response.Code, response.Body.String())
	}
}

func TestGameEndpointsRejectInvalidQueryAndMissingGame(t *testing.T) {
	router, _ := testRouter(t, "http://localhost:3000", 10)
	tests := []struct {
		path string
		want int
	}{
		{"/v1/games?page=0", http.StatusBadRequest},
		{"/v1/games?pageSize=51", http.StatusBadRequest},
		{"/v1/games?minRating=invalid", http.StatusBadRequest},
		{"/v1/games?sort=invalid", http.StatusBadRequest},
		{"/v1/games/missing", http.StatusNotFound},
	}
	for _, test := range tests {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, test.path, nil))
		if response.Code != test.want || !strings.Contains(response.Body.String(), `"requestId"`) {
			t.Fatalf("%s: status=%d body=%s", test.path, response.Code, response.Body.String())
		}
	}
}

func TestCreateReviewValidatesAndReturnsCreatedReview(t *testing.T) {
	router, _ := testRouter(t, "http://localhost:3000", 10)
	request := httptest.NewRequest(http.MethodPost, "/v1/games/game-002/reviews", strings.NewReader(`{"reviewerName":"Ada","rating":5,"text":"A beautifully focused game with excellent pacing."}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusCreated || !strings.Contains(response.Body.String(), `"reviewerName":"Ada"`) || response.Header().Get("Location") == "" {
		t.Fatalf("unexpected create response: status=%d headers=%v body=%s", response.Code, response.Header(), response.Body.String())
	}

	request = httptest.NewRequest(http.MethodPost, "/v1/games/game-002/reviews", strings.NewReader(`{"reviewerName":"A","rating":8,"text":"short"}`))
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusUnprocessableEntity || !strings.Contains(response.Body.String(), `"fields"`) {
		t.Fatalf("unexpected validation response: status=%d body=%s", response.Code, response.Body.String())
	}
}

func TestCreateReviewRejectsUnsafeTransportInput(t *testing.T) {
	router, _ := testRouter(t, "http://localhost:3000", 10)
	tests := []struct {
		contentType string
		body        string
		want        int
	}{
		{"text/plain", `{}`, http.StatusUnsupportedMediaType},
		{"application/json", `{"reviewerName":"Ada","rating":5,"text":"A long enough review.","admin":true}`, http.StatusBadRequest},
		{"application/json", `{invalid`, http.StatusBadRequest},
		{"application/json", string(bytes.Repeat([]byte("a"), maximumReviewBody+1)), http.StatusBadRequest},
	}
	for _, test := range tests {
		request := httptest.NewRequest(http.MethodPost, "/v1/games/game-002/reviews", strings.NewReader(test.body))
		request.Header.Set("Content-Type", test.contentType)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != test.want {
			t.Fatalf("contentType=%s: status=%d body=%s", test.contentType, response.Code, response.Body.String())
		}
	}
}

func TestCreateReviewRateLimit(t *testing.T) {
	router, _ := testRouter(t, "http://localhost:3000", 1)
	for attempt := 0; attempt < 2; attempt++ {
		request := httptest.NewRequest(http.MethodPost, "/v1/games/game-002/reviews", strings.NewReader(`{"reviewerName":"Ada","rating":5,"text":"A beautifully focused game with excellent pacing."}`))
		request.Header.Set("Content-Type", "application/json")
		request.RemoteAddr = "192.0.2.1:1234"
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if attempt == 1 && (response.Code != http.StatusTooManyRequests || response.Header().Get("Retry-After") == "") {
			t.Fatalf("rate limit response: status=%d headers=%v body=%s", response.Code, response.Header(), response.Body.String())
		}
	}
}

func TestCORSAndSecurityHeaders(t *testing.T) {
	router, _ := testRouter(t, "http://localhost:3000", 10)
	request := httptest.NewRequest(http.MethodOptions, "/v1/games", nil)
	request.Header.Set("Origin", "http://localhost:3000")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent || response.Header().Get("Access-Control-Allow-Origin") == "" || response.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("unexpected preflight response: status=%d headers=%v", response.Code, response.Header())
	}
}

func TestWebSocketReceivesCreatedReview(t *testing.T) {
	metrics := &observability.Metrics{}
	var router http.Handler
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) { router.ServeHTTP(writer, request) }))
	defer server.Close()
	games, reviews, _ := memory.Seed()
	store, _ := memory.New(games, reviews)
	hub := realtime.NewHub([]string{server.URL}, 10, metrics)
	router = NewRouter(Settings{
		Games: gameservices.New(store, store), Reviews: reviewservices.New(store, store, hub, reviewservices.SystemClock{}, reviewservices.RandomIDGenerator{}),
		WebSocket: hub, Metrics: metrics, Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		AllowedOrigins: []string{server.URL}, ReviewRateLimit: 10, ReviewRateWindow: time.Minute,
	})
	websocketURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/v1/ws"
	connection, _, err := websocket.Dial(context.Background(), websocketURL, &websocket.DialOptions{HTTPHeader: http.Header{"Origin": []string{server.URL}}})
	if err != nil {
		t.Fatal(err)
	}
	defer connection.CloseNow()

	request, _ := http.NewRequest(http.MethodPost, server.URL+"/v1/games/game-002/reviews", strings.NewReader(`{"reviewerName":"Live Reviewer","rating":5,"text":"This review should be delivered over the live connection."}`))
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(response.Body)
		t.Fatalf("create status=%d body=%s", response.StatusCode, body)
	}

	readContext, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, payload, err := connection.Read(readContext)
	if err != nil {
		t.Fatal(err)
	}
	var event realtime.ReviewCreatedEvent
	if err = json.Unmarshal(payload, &event); err != nil {
		t.Fatal(err)
	}
	if event.Type != "review.created" || event.GameID != "game-002" || event.Review.ReviewerName != "Live Reviewer" {
		t.Fatalf("unexpected event: %#v", event)
	}
}
