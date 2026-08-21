package ports

import (
	"context"

	"github.com/septivan/viger/backend/internal/core/game/domain"
)

type Repository interface {
	List(context.Context) ([]domain.Game, error)
	FindByID(context.Context, string) (domain.Game, bool, error)
}

