package memory

import (
	"context"
	"fmt"
	"sync"

	game "github.com/septivan/viger/backend/internal/core/game/domain"
	review "github.com/septivan/viger/backend/internal/core/review/domain"
)

type Store struct {
	mutex   sync.RWMutex
	games   map[string]game.Game
	order   []string
	reviews map[string][]review.Review
}

func New(games []game.Game, reviews []review.Review) (*Store, error) {
	store := &Store{
		games:   make(map[string]game.Game, len(games)),
		order:   make([]string, 0, len(games)),
		reviews: make(map[string][]review.Review),
	}
	for _, item := range games {
		if _, exists := store.games[item.ID]; exists {
			return nil, fmt.Errorf("duplicate game ID %q", item.ID)
		}
		store.games[item.ID] = cloneGame(item)
		store.order = append(store.order, item.ID)
	}
	for _, item := range reviews {
		if _, exists := store.games[item.GameID]; !exists {
			return nil, fmt.Errorf("review %q references missing game %q", item.ID, item.GameID)
		}
		store.reviews[item.GameID] = append(store.reviews[item.GameID], item)
	}
	return store, nil
}

func (store *Store) List(_ context.Context) ([]game.Game, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	items := make([]game.Game, 0, len(store.order))
	for _, id := range store.order {
		items = append(items, cloneGame(store.games[id]))
	}
	return items, nil
}

func (store *Store) FindByID(_ context.Context, id string) (game.Game, bool, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	item, found := store.games[id]
	return cloneGame(item), found, nil
}

func (store *Store) ListByGame(_ context.Context, gameID string) ([]review.Review, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	items := store.reviews[gameID]
	result := make([]review.Review, len(items))
	copy(result, items)
	return result, nil
}

func (store *Store) Create(_ context.Context, item review.Review) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if _, found := store.games[item.GameID]; !found {
		return fmt.Errorf("game %q does not exist", item.GameID)
	}
	store.reviews[item.GameID] = append(store.reviews[item.GameID], item)
	return nil
}

func cloneGame(item game.Game) game.Game {
	item.Platforms = append([]string(nil), item.Platforms...)
	return item
}
