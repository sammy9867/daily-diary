package repository

import (
	"github.com/sammy9867/daily-diary/backend/domain"
)

// AuthRepository handles users authentication
type AuthRepository interface {
	Login(email, password string) (*domain.TokenDetail, error)
	Logout(uuid string) (int64, error)
}
