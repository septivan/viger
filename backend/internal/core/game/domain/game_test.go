package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewGameValidatesRequiredFields(t *testing.T) {
	valid := NewGameInput{ID: "game-1", Title: "Hades", Description: "A fast roguelike adventure.", Genre: "Roguelike", Platforms: []string{"PC"}, Developer: "Studio", ReleaseDate: time.Now()}
	tests := []struct {
		name   string
		mutate func(*NewGameInput)
		want   error
	}{
		{"ID", func(input *NewGameInput) { input.ID = "" }, ErrInvalidID},
		{"title", func(input *NewGameInput) { input.Title = "" }, ErrInvalidTitle},
		{"description", func(input *NewGameInput) { input.Description = "short" }, ErrInvalidDescription},
		{"genre", func(input *NewGameInput) { input.Genre = "" }, ErrInvalidGenre},
		{"platform", func(input *NewGameInput) { input.Platforms = nil }, ErrInvalidPlatforms},
		{"developer", func(input *NewGameInput) { input.Developer = "" }, ErrInvalidDeveloper},
		{"release date", func(input *NewGameInput) { input.ReleaseDate = time.Time{} }, ErrInvalidReleaseDate},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := valid
			test.mutate(&input)
			_, err := New(input)
			if !errors.Is(err, test.want) {
				t.Fatalf("error = %v, want %v", err, test.want)
			}
		})
	}
}

