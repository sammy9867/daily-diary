package usecase

import (
	"github.com/sammy9867/daily-diary/backend/domain"
)

// AuthUseCase represents the users authentication usecase
type AuthUseCase interface {
	Login(email, password string) (*domain.TokenDetail, error)
	Logout(uuid string) (int64, error)
}
