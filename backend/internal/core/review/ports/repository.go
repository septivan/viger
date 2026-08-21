package ports

import (
	"context"

	"github.com/septivan/viger/backend/internal/core/review/domain"
)

// Repository is the storage boundary required by review use cases.
type Repository interface {
	ListByGame(context.Context, string) ([]domain.Review, error)
	Create(context.Context, domain.Review) error
}

// EventPublisher announces a review only after it has been stored successfully.
type EventPublisher interface {
	PublishReviewCreated(context.Context, domain.Review) error
}
