package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewReviewValidatesInput(t *testing.T) {
	valid := NewReviewInput{ID: "review-1", GameID: "game-1", ReviewerName: "Alex", Rating: 5, Text: "A thoughtful and polished experience.", CreatedAt: time.Now()}
	tests := []struct {
		name   string
		mutate func(*NewReviewInput)
		want   error
	}{
		{"missing ID", func(input *NewReviewInput) { input.ID = "" }, ErrInvalidID},
		{"missing game", func(input *NewReviewInput) { input.GameID = "" }, ErrInvalidGameID},
		{"short name", func(input *NewReviewInput) { input.ReviewerName = "A" }, ErrInvalidReviewerName},
		{"long name", func(input *NewReviewInput) { input.ReviewerName = strings.Repeat("a", 81) }, ErrInvalidReviewerName},
		{"low rating", func(input *NewReviewInput) { input.Rating = 0 }, ErrInvalidRating},
		{"high rating", func(input *NewReviewInput) { input.Rating = 6 }, ErrInvalidRating},
		{"short text", func(input *NewReviewInput) { input.Text = "too short" }, ErrInvalidText},
		{"control character", func(input *NewReviewInput) { input.Text = "Valid length\x00 invalid" }, ErrControlCharacters},
		{"missing time", func(input *NewReviewInput) { input.CreatedAt = time.Time{} }, ErrInvalidCreatedAt},
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

func TestNewReviewTrimsAndNormalizesTime(t *testing.T) {
	createdAt := time.Date(2026, 8, 21, 12, 0, 0, 0, time.FixedZone("local", 7*60*60))
	review, err := New(NewReviewInput{
		ID: " review-1 ", GameID: " game-1 ", ReviewerName: " Alex ", Rating: 4,
		Text: " A detailed review with useful context. ", CreatedAt: createdAt,
	})
	if err != nil {
		t.Fatal(err)
	}
	if review.ID != "review-1" || review.ReviewerName != "Alex" || review.Text[0] == ' ' {
		t.Fatalf("review was not normalized: %#v", review)
	}
	if review.CreatedAt.Location() != time.UTC {
		t.Fatalf("created time location = %v, want UTC", review.CreatedAt.Location())
	}
}
