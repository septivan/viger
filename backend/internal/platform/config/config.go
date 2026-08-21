package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddress          string
	AllowedOrigins       []string
	ReviewRateLimit      int
	ReviewRateWindow     time.Duration
	MaximumWSConnections int
}

func Load() (Config, error) {
	rateLimit, err := positiveInteger("VIGER_REVIEW_RATE_LIMIT", 10)
	if err != nil {
		return Config{}, err
	}
	maximumConnections, err := positiveInteger("VIGER_MAX_WS_CONNECTIONS", 500)
	if err != nil {
		return Config{}, err
	}
	origins := splitValues(environment("VIGER_ALLOWED_ORIGINS", "http://localhost:3000"))
	if len(origins) == 0 {
		return Config{}, fmt.Errorf("VIGER_ALLOWED_ORIGINS must contain at least one origin")
	}
	return Config{
		HTTPAddress:          environment("VIGER_HTTP_ADDRESS", ":8080"),
		AllowedOrigins:       origins,
		ReviewRateLimit:      rateLimit,
		ReviewRateWindow:     time.Minute,
		MaximumWSConnections: maximumConnections,
	}, nil
}

func environment(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

func positiveInteger(name string, fallback int) (int, error) {
	value := environment(name, strconv.Itoa(fallback))
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}
	return parsed, nil
}

func splitValues(value string) []string {
	result := make([]string, 0)
	for _, item := range strings.Split(value, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, strings.TrimSuffix(trimmed, "/"))
		}
	}
	return result
}
