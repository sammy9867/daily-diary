package model

import (
	"time"
)

// Entry model struct
type Entry struct {
	ID          uint64       `gorm:"primary_key;auto_increment" json:"id"`
	Title       string       `gorm:"size:255;not null" json:"title"`
	Description string       `gorm:"size:100;not null" json:"description"`
	EntryImages []EntryImage `gorm:"foreignkey:entry_id" json:"images"`
	Owner       User         `json:"owner"`
	OwnerID     uint64       `gorm:"not null" json:"owner_id"`
	CreatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
