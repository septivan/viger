package ports

import (
	"context"

	"github.com/septivan/viger/backend/internal/core/review/domain"
)

type Repository interface {
	ListByGame(context.Context, string) ([]domain.Review, error)
	Create(context.Context, domain.Review) error
}

type EventPublisher interface {
	PublishReviewCreated(context.Context, domain.Review) error
}
