package services

import (
	"context"
	"errors"
	"testing"

	"github.com/septivan/viger/backend/internal/adapters/outbound/memory"
)

func testService(t *testing.T) Service {
	t.Helper()
	games, reviews, err := memory.Seed()
	if err != nil {
		t.Fatal(err)
	}
	store, err := memory.New(games, reviews)
	if err != nil {
		t.Fatal(err)
	}
	return New(store, store)
}

func TestListSearchesFiltersSortsAndPaginates(t *testing.T) {
	service := testService(t)
	page, err := service.List(context.Background(), ListQuery{Search: "hades", Genre: "Roguelike", Platform: "PC", Sort: "title_asc", Page: 1, PageSize: 12})
	if err != nil {
		t.Fatal(err)
	}
	if page.TotalItems != 1 || len(page.Items) != 1 || page.Items[0].Game.Title != "Hades" {
		t.Fatalf("unexpected page: %#v", page)
	}

	page, err = service.List(context.Background(), ListQuery{Sort: "newest", Page: 2, PageSize: 12})
	if err != nil {
		t.Fatal(err)
	}
	if page.TotalItems != 48 || page.TotalPages != 4 || len(page.Items) != 12 {
		t.Fatalf("unexpected pagination: %#v", page)
	}
}

func TestListRejectsInvalidQueries(t *testing.T) {
	service := testService(t)
	tests := []struct {
		query ListQuery
		want  error
	}{
		{ListQuery{Page: 0, PageSize: 12}, ErrInvalidPage},
		{ListQuery{Page: 1, PageSize: 51}, ErrInvalidPageSize},
		{ListQuery{Page: 1, PageSize: 12, MinRating: 6}, ErrInvalidMinRating},
		{ListQuery{Page: 1, PageSize: 12, Sort: "unknown"}, ErrInvalidGameSort},
	}
	for _, test := range tests {
		if _, err := service.List(context.Background(), test.query); !errors.Is(err, test.want) {
			t.Fatalf("error = %v, want %v", err, test.want)
		}
	}
}

func TestListReviewsReturnsNewestFirst(t *testing.T) {
	service := testService(t)
	page, err := service.ListReviews(context.Background(), "game-002", ReviewQuery{Sort: "newest", Page: 1, PageSize: 5})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) < 2 || page.Items[0].CreatedAt.Before(page.Items[1].CreatedAt) {
		t.Fatalf("reviews are not newest first: %#v", page.Items)
	}
}

