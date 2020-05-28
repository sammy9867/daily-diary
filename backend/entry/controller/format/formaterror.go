package format

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/sammy9867/daily-diary/backend/domain"
)

// Initialize is used to initialize the entry
func Initialize(e *domain.Entry) {
	e.ID = 0
	e.Title = html.EscapeString(strings.TrimSpace(e.Title))
	e.Description = html.EscapeString(strings.TrimSpace(e.Description))
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
}

// Validate is used to check if the entry has a correct input format
func Validate(e *domain.Entry) error {
	if e.Title == "" {
		return errors.New("Required Title")
	}
	if e.Description == "" {
		return errors.New("Required Description")
	}
	if e.OwnerID < 1 {
		return errors.New("Required Owner of the Post")
	}
	return nil
}
