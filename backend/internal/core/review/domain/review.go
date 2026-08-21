package domain

import (
	"errors"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidID           = errors.New("review ID is required")
	ErrInvalidGameID       = errors.New("game ID is required")
	ErrInvalidReviewerName = errors.New("reviewer name must contain between 2 and 80 characters")
	ErrInvalidRating       = errors.New("rating must be an integer between 1 and 5")
	ErrInvalidText         = errors.New("review text must contain between 10 and 2000 characters")
	ErrInvalidCreatedAt    = errors.New("review creation time is required")
	ErrControlCharacters   = errors.New("review fields contain unsupported control characters")
)

type Review struct {
	ID           string
	GameID       string
	ReviewerName string
	Rating       int
	Text         string
	CreatedAt    time.Time
}

type NewReviewInput struct {
	ID           string
	GameID       string
	ReviewerName string
	Rating       int
	Text         string
	CreatedAt    time.Time
}

func New(input NewReviewInput) (Review, error) {
	input.ID = strings.TrimSpace(input.ID)
	input.GameID = strings.TrimSpace(input.GameID)
	input.ReviewerName = strings.TrimSpace(input.ReviewerName)
	input.Text = strings.TrimSpace(input.Text)

	if input.ID == "" {
		return Review{}, ErrInvalidID
	}
	if input.GameID == "" {
		return Review{}, ErrInvalidGameID
	}
	if length := utf8.RuneCountInString(input.ReviewerName); length < 2 || length > 80 {
		return Review{}, ErrInvalidReviewerName
	}
	if input.Rating < 1 || input.Rating > 5 {
		return Review{}, ErrInvalidRating
	}
	if length := utf8.RuneCountInString(input.Text); length < 10 || length > 2000 {
		return Review{}, ErrInvalidText
	}
	if containsUnsupportedControl(input.ReviewerName, false) || containsUnsupportedControl(input.Text, true) {
		return Review{}, ErrControlCharacters
	}
	if input.CreatedAt.IsZero() {
		return Review{}, ErrInvalidCreatedAt
	}

	return Review{
		ID:           input.ID,
		GameID:       input.GameID,
		ReviewerName: input.ReviewerName,
		Rating:       input.Rating,
		Text:         input.Text,
		CreatedAt:    input.CreatedAt.UTC(),
	}, nil
}

func containsUnsupportedControl(value string, allowWhitespace bool) bool {
	for _, character := range value {
		if !unicode.IsControl(character) {
			continue
		}
		if allowWhitespace && (character == '\n' || character == '\t' || character == '\r') {
			continue
		}
		return true
	}
	return false
}
