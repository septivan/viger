package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/septivan/viger/backend/internal/adapters/outbound/memory"
	review "github.com/septivan/viger/backend/internal/core/review/domain"
)

type fixedClock struct{ value time.Time }

func (clock fixedClock) Now() time.Time { return clock.value }

type fixedID struct{ value string }

func (generator fixedID) NewID() (string, error) { return generator.value, nil }

type eventRecorder struct{ created []review.Review }

func (recorder *eventRecorder) PublishReviewCreated(_ context.Context, item review.Review) error {
	recorder.created = append(recorder.created, item)
	return nil
}

func TestCreatePersistsAndPublishesValidatedReview(t *testing.T) {
	games, reviews, err := memory.Seed()
	if err != nil {
		t.Fatal(err)
	}
	store, err := memory.New(games, reviews)
	if err != nil {
		t.Fatal(err)
	}
	events := &eventRecorder{}
	createdAt := time.Date(2026, 8, 21, 12, 30, 0, 0, time.UTC)
	service := New(store, store, events, fixedClock{createdAt}, fixedID{"review-new"})

	created, err := service.Create(context.Background(), CreateInput{GameID: games[0].ID, ReviewerName: " Alex ", Rating: 5, Text: " A polished and memorable experience. "})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID != "review-new" || created.ReviewerName != "Alex" || !created.CreatedAt.Equal(createdAt) {
		t.Fatalf("unexpected review: %#v", created)
	}
	stored, err := store.ListByGame(context.Background(), games[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored[len(stored)-1].ID != created.ID || len(events.created) != 1 || events.created[0].ID != created.ID {
		t.Fatalf("review was not stored and published: stored=%#v events=%#v", stored, events.created)
	}
}

func TestCreateRejectsMissingGameAndInvalidReview(t *testing.T) {
	games, reviews, _ := memory.Seed()
	store, _ := memory.New(games, reviews)
	service := New(store, store, &eventRecorder{}, fixedClock{time.Now()}, fixedID{"review-new"})
	if _, err := service.Create(context.Background(), CreateInput{GameID: "missing", ReviewerName: "Alex", Rating: 5, Text: "A valid review body for a missing game."}); !errors.Is(err, ErrGameNotFound) {
		t.Fatalf("error = %v, want %v", err, ErrGameNotFound)
	}
	if _, err := service.Create(context.Background(), CreateInput{GameID: games[0].ID, ReviewerName: "A", Rating: 9, Text: "short"}); err == nil {
		t.Fatal("invalid review was accepted")
	}
}
