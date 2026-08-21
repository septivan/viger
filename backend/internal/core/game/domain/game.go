package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrInvalidID          = errors.New("game ID is required")
	ErrInvalidTitle       = errors.New("game title must contain between 1 and 120 characters")
	ErrInvalidDescription = errors.New("game description must contain between 10 and 1000 characters")
	ErrInvalidGenre       = errors.New("game genre is required")
	ErrInvalidPlatforms   = errors.New("game must have at least one platform")
	ErrInvalidDeveloper   = errors.New("game developer is required")
	ErrInvalidReleaseDate = errors.New("game release date is required")
)

type Game struct {
	ID          string
	Title       string
	Description string
	Genre       string
	Platforms   []string
	Developer   string
	ReleaseDate time.Time
}

type NewGameInput struct {
	ID          string
	Title       string
	Description string
	Genre       string
	Platforms   []string
	Developer   string
	ReleaseDate time.Time
}

func New(input NewGameInput) (Game, error) {
	input.ID = strings.TrimSpace(input.ID)
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	input.Genre = strings.TrimSpace(input.Genre)
	input.Developer = strings.TrimSpace(input.Developer)

	if input.ID == "" {
		return Game{}, ErrInvalidID
	}
	if length := utf8.RuneCountInString(input.Title); length < 1 || length > 120 {
		return Game{}, ErrInvalidTitle
	}
	if length := utf8.RuneCountInString(input.Description); length < 10 || length > 1000 {
		return Game{}, ErrInvalidDescription
	}
	if input.Genre == "" {
		return Game{}, ErrInvalidGenre
	}
	if len(input.Platforms) == 0 {
		return Game{}, ErrInvalidPlatforms
	}
	platforms := make([]string, 0, len(input.Platforms))
	for _, platform := range input.Platforms {
		if trimmed := strings.TrimSpace(platform); trimmed != "" {
			platforms = append(platforms, trimmed)
		}
	}
	if len(platforms) == 0 {
		return Game{}, ErrInvalidPlatforms
	}
	if input.Developer == "" {
		return Game{}, ErrInvalidDeveloper
	}
	if input.ReleaseDate.IsZero() {
		return Game{}, ErrInvalidReleaseDate
	}

	return Game{
		ID:          input.ID,
		Title:       input.Title,
		Description: input.Description,
		Genre:       input.Genre,
		Platforms:   platforms,
		Developer:   input.Developer,
		ReleaseDate: input.ReleaseDate.UTC(),
	}, nil
}

