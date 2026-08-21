package memory

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	review "github.com/septivan/viger/backend/internal/core/review/domain"
)

func TestSeedIsDeterministicAndSubstantial(t *testing.T) {
	firstGames, firstReviews, err := Seed()
	if err != nil {
		t.Fatal(err)
	}
	secondGames, secondReviews, err := Seed()
	if err != nil {
		t.Fatal(err)
	}
	if len(firstGames) != 48 {
		t.Fatalf("games = %d, want 48", len(firstGames))
	}
	if len(firstReviews) < 250 || len(firstReviews) > 400 {
		t.Fatalf("reviews = %d, want between 250 and 400", len(firstReviews))
	}
	if firstGames[10].ID != secondGames[10].ID || firstGames[10].Title != secondGames[10].Title || firstReviews[20] != secondReviews[20] {
		t.Fatal("seed data is not deterministic")
	}
}

func TestStoreSupportsConcurrentReviewCreation(t *testing.T) {
	games, reviews, err := Seed()
	if err != nil {
		t.Fatal(err)
	}
	store, err := New(games, reviews)
	if err != nil {
		t.Fatal(err)
	}
	const additions = 40
	var wait sync.WaitGroup
	for index := 0; index < additions; index++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			item, createErr := review.New(review.NewReviewInput{
				ID: fmt.Sprintf("concurrent-%d", index), GameID: games[0].ID,
				ReviewerName: "Concurrent Reviewer", Rating: 4,
				Text: "A review created during the concurrency test.", CreatedAt: time.Now(),
			})
			if createErr != nil {
				t.Error(createErr)
				return
			}
			if createErr = store.Create(context.Background(), item); createErr != nil {
				t.Error(createErr)
			}
		}(index)
	}
	wait.Wait()
	items, err := store.ListByGame(context.Background(), games[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	initial := 0
	for _, item := range reviews {
		if item.GameID == games[0].ID {
			initial++
		}
	}
	if len(items) != initial+additions {
		t.Fatalf("reviews = %d, want %d", len(items), initial+additions)
	}
}
