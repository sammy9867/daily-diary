package domain

import (
	"time"
)

// EntryImage model struct
type EntryImage struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	URL       string    `gorm:"size:255;not null" json:"image_url"`
	EntryID   uint64    `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
