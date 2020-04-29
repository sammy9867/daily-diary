package usecase

import (
	"github.com/sammy9867/daily-diary/backend/user/model"
)

// UserUseCase represents the users usecase
type UserUseCase interface {
	CreateUser(*model.User) (*model.User, error)
	UpdateUser(uint64, *model.User) (*model.User, error)
	DeleteUser(uint64) (int64, error)
	GetUserByID(uint64) (*model.User, error)
	GetAllUsers() (*[]model.User, error)
}
