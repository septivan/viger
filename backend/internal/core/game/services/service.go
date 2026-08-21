package services

import (
	"context"
	"errors"
	"math"
	"sort"
	"strings"

	"github.com/septivan/viger/backend/internal/core/game/domain"
	gameports "github.com/septivan/viger/backend/internal/core/game/ports"
	reviewdomain "github.com/septivan/viger/backend/internal/core/review/domain"
	reviewport "github.com/septivan/viger/backend/internal/core/review/ports"
)

var (
	ErrGameNotFound      = errors.New("game not found")
	ErrInvalidPage       = errors.New("page must be at least 1")
	ErrInvalidPageSize   = errors.New("page size must be between 1 and 50")
	ErrInvalidSearch     = errors.New("search query must not exceed 100 characters")
	ErrInvalidMinRating  = errors.New("minimum rating must be between 0 and 5")
	ErrInvalidGameSort   = errors.New("unsupported game sort")
	ErrInvalidReviewSort = errors.New("unsupported review sort")
)

// Service owns catalog and game-detail read use cases.
type Service struct {
	games   gameports.Repository
	reviews reviewport.Repository
}

// New creates a catalog service backed only by core repository ports.
func New(games gameports.Repository, reviews reviewport.Repository) Service {
	return Service{games: games, reviews: reviews}
}

type ListQuery struct {
	Search    string
	Genre     string
	Platform  string
	MinRating float64
	Sort      string
	Page      int
	PageSize  int
}

type RatingSummary struct {
	Average      float64
	Total        int
	Distribution map[int]int
}

type CatalogItem struct {
	Game    domain.Game
	Ratings RatingSummary
}

type Page[T any] struct {
	Items      []T
	Page       int
	PageSize   int
	TotalItems int
	TotalPages int
}

func (service Service) List(context context.Context, query ListQuery) (Page[CatalogItem], error) {
	if err := validatePage(query.Page, query.PageSize); err != nil {
		return Page[CatalogItem]{}, err
	}
	query.Search = strings.TrimSpace(query.Search)
	if len([]rune(query.Search)) > 100 {
		return Page[CatalogItem]{}, ErrInvalidSearch
	}
	if query.MinRating < 0 || query.MinRating > 5 {
		return Page[CatalogItem]{}, ErrInvalidMinRating
	}
	if query.Sort == "" {
		query.Sort = "rating_desc"
	}
	if !validGameSort(query.Sort) {
		return Page[CatalogItem]{}, ErrInvalidGameSort
	}

	games, err := service.games.List(context)
	if err != nil {
		return Page[CatalogItem]{}, err
	}
	items := make([]CatalogItem, 0, len(games))
	for _, game := range games {
		reviews, listErr := service.reviews.ListByGame(context, game.ID)
		if listErr != nil {
			return Page[CatalogItem]{}, listErr
		}
		item := CatalogItem{Game: game, Ratings: summarize(reviews)}
		if matches(item, query) {
			items = append(items, item)
		}
	}
	sortGames(items, query.Sort)
	return paginate(items, query.Page, query.PageSize), nil
}

func (service Service) Find(context context.Context, gameID string) (CatalogItem, error) {
	game, found, err := service.games.FindByID(context, strings.TrimSpace(gameID))
	if err != nil {
		return CatalogItem{}, err
	}
	if !found {
		return CatalogItem{}, ErrGameNotFound
	}
	reviews, err := service.reviews.ListByGame(context, game.ID)
	if err != nil {
		return CatalogItem{}, err
	}
	return CatalogItem{Game: game, Ratings: summarize(reviews)}, nil
}

type ReviewQuery struct {
	Sort     string
	Page     int
	PageSize int
}

func (service Service) ListReviews(context context.Context, gameID string, query ReviewQuery) (Page[reviewdomain.Review], error) {
	if _, err := service.Find(context, gameID); err != nil {
		return Page[reviewdomain.Review]{}, err
	}
	if err := validatePage(query.Page, query.PageSize); err != nil {
		return Page[reviewdomain.Review]{}, err
	}
	if query.Sort == "" {
		query.Sort = "newest"
	}
	if !validReviewSort(query.Sort) {
		return Page[reviewdomain.Review]{}, ErrInvalidReviewSort
	}
	reviews, err := service.reviews.ListByGame(context, gameID)
	if err != nil {
		return Page[reviewdomain.Review]{}, err
	}
	sortReviews(reviews, query.Sort)
	return paginate(reviews, query.Page, query.PageSize), nil
}

func validatePage(page, pageSize int) error {
	if page < 1 {
		return ErrInvalidPage
	}
	if pageSize < 1 || pageSize > 50 {
		return ErrInvalidPageSize
	}
	return nil
}

func summarize(reviews []reviewdomain.Review) RatingSummary {
	summary := RatingSummary{Total: len(reviews), Distribution: map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}}
	if len(reviews) == 0 {
		return summary
	}
	total := 0
	for _, review := range reviews {
		total += review.Rating
		summary.Distribution[review.Rating]++
	}
	summary.Average = math.Round((float64(total)/float64(len(reviews)))*10) / 10
	return summary
}

func matches(item CatalogItem, query ListQuery) bool {
	if query.Search != "" && !strings.Contains(strings.ToLower(item.Game.Title), strings.ToLower(query.Search)) {
		return false
	}
	if query.Genre != "" && !strings.EqualFold(item.Game.Genre, query.Genre) {
		return false
	}
	if query.Platform != "" {
		found := false
		for _, platform := range item.Game.Platforms {
			found = found || strings.EqualFold(platform, query.Platform)
		}
		if !found {
			return false
		}
	}
	return item.Ratings.Average >= query.MinRating
}

func validGameSort(value string) bool {
	return value == "rating_desc" || value == "reviews_desc" || value == "title_asc" || value == "newest"
}

func sortGames(items []CatalogItem, ordering string) {
	sort.SliceStable(items, func(left, right int) bool {
		switch ordering {
		case "reviews_desc":
			return items[left].Ratings.Total > items[right].Ratings.Total
		case "title_asc":
			return strings.ToLower(items[left].Game.Title) < strings.ToLower(items[right].Game.Title)
		case "newest":
			return items[left].Game.ReleaseDate.After(items[right].Game.ReleaseDate)
		default:
			if items[left].Ratings.Average == items[right].Ratings.Average {
				return items[left].Ratings.Total > items[right].Ratings.Total
			}
			return items[left].Ratings.Average > items[right].Ratings.Average
		}
	})
}

func validReviewSort(value string) bool {
	return value == "newest" || value == "oldest" || value == "rating_desc" || value == "rating_asc"
}

func sortReviews(reviews []reviewdomain.Review, ordering string) {
	sort.SliceStable(reviews, func(left, right int) bool {
		switch ordering {
		case "oldest":
			return reviews[left].CreatedAt.Before(reviews[right].CreatedAt)
		case "rating_desc":
			return reviews[left].Rating > reviews[right].Rating
		case "rating_asc":
			return reviews[left].Rating < reviews[right].Rating
		default:
			return reviews[left].CreatedAt.After(reviews[right].CreatedAt)
		}
	})
}

func paginate[T any](items []T, page, pageSize int) Page[T] {
	total := len(items)
	totalPages := 0
	if total > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}
	start := (page - 1) * pageSize
	if start >= total {
		return Page[T]{Items: []T{}, Page: page, PageSize: pageSize, TotalItems: total, TotalPages: totalPages}
	}
	end := min(start+pageSize, total)
	return Page[T]{Items: items[start:end], Page: page, PageSize: pageSize, TotalItems: total, TotalPages: totalPages}
}
