package ports

import (
	"context"

	"github.com/septivan/viger/backend/internal/core/game/domain"
)

// Repository is the storage boundary required by game catalog use cases.
type Repository interface {
	List(context.Context) ([]domain.Game, error)
	FindByID(context.Context, string) (domain.Game, bool, error)
}
