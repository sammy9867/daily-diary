package usecase

import (
	"github.com/sammy9867/daily-diary/backend/model"
)

// UserUseCase represents the users usecase
type UserUseCase interface {
	SignIn(email, password string) (string, error)
	CreateUser(*model.User) (*model.User, error)
	UpdateUser(uint64, *model.User) (*model.User, error)
	DeleteUser(uint64) (int64, error)
	GetUserByID(uint64) (*model.User, error)
	GetAllUsers() (*[]model.User, error)
}
