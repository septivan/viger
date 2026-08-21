package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	gameports "github.com/septivan/viger/backend/internal/core/game/ports"
	"github.com/septivan/viger/backend/internal/core/review/domain"
	"github.com/septivan/viger/backend/internal/core/review/ports"
)

var ErrGameNotFound = errors.New("game not found")

// Clock keeps review creation deterministic in tests.
type Clock interface {
	Now() time.Time
}

// IDGenerator supplies storage-independent review identifiers.
type IDGenerator interface {
	NewID() (string, error)
}

// Service owns the review creation use case.
type Service struct {
	games   gameports.Repository
	reviews ports.Repository
	events  ports.EventPublisher
	clock   Clock
	ids     IDGenerator
}

// New wires review creation to its required core ports.
func New(games gameports.Repository, reviews ports.Repository, events ports.EventPublisher, clock Clock, ids IDGenerator) Service {
	return Service{games: games, reviews: reviews, events: events, clock: clock, ids: ids}
}

type CreateInput struct {
	GameID       string
	ReviewerName string
	Rating       int
	Text         string
}

// Create validates the game relationship, stores the review, then publishes it.
func (service Service) Create(context context.Context, input CreateInput) (domain.Review, error) {
	if _, found, err := service.games.FindByID(context, input.GameID); err != nil {
		return domain.Review{}, err
	} else if !found {
		return domain.Review{}, ErrGameNotFound
	}
	id, err := service.ids.NewID()
	if err != nil {
		return domain.Review{}, fmt.Errorf("generate review ID: %w", err)
	}
	review, err := domain.New(domain.NewReviewInput{
		ID: id, GameID: input.GameID, ReviewerName: input.ReviewerName,
		Rating: input.Rating, Text: input.Text, CreatedAt: service.clock.Now(),
	})
	if err != nil {
		return domain.Review{}, err
	}
	if err = service.reviews.Create(context, review); err != nil {
		return domain.Review{}, fmt.Errorf("store review: %w", err)
	}
	if service.events != nil {
		if err = service.events.PublishReviewCreated(context, review); err != nil {
			return review, fmt.Errorf("publish review event: %w", err)
		}
	}
	return review, nil
}

type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now().UTC() }

type RandomIDGenerator struct{}

func (RandomIDGenerator) NewID() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return "review-" + hex.EncodeToString(value), nil
}
