package domain

import (
	"html"
	"strings"
	"time"
)

// User model struct
type User struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id" redis:"id"`
	Username  string    `gorm:"size:255;not null;unique" json:"username" redis:"username"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"  redis:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"  redis:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at" redis:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at" redis:"updated_at"`
}

// Initialize is used to initialize the user before creating/login/updating
func (u *User) Initialize() {
	u.ID = 0
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}
