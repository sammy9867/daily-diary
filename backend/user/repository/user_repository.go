package repository

import (
	"github.com/sammy9867/daily-diary/backend/user/model"
)

// UserRepository represents the users repository
type UserRepository interface {
	SignIn(email, password string) (string, error)
	CreateUser(*model.User) (*model.User, error)
	UpdateUser(uint64, *model.User) (*model.User, error)
	DeleteUser(uid uint64) (int64, error)
	GetUserByID(uid uint64) (*model.User, error)
	GetAllUsers() (*[]model.User, error)
}
