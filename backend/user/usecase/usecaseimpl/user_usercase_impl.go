package usecaseimpl

import (
	"github.com/sammy9867/daily-diary/backend/user/model"
	"github.com/sammy9867/daily-diary/backend/user/repository"
	"github.com/sammy9867/daily-diary/backend/user/usecase"
)

type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserUseCase does
func NewUserUseCase(ur repository.UserRepository) usecase.UserUseCase {
	return &userUsecase{userRepo: ur}
}

func (userUC *userUsecase) SignIn(email, password string) (string, error) {
	token, err := userUC.userRepo.SignIn(email, password)
	if err != nil {
		return "error signing in user", nil
	}

	return token, nil
}

func (userUC *userUsecase) CreateUser(u *model.User) (*model.User, error) {
	createdUser, err := userUC.userRepo.CreateUser(u)
	if err != nil {
		return &model.User{}, nil
	}

	return createdUser, nil
}

func (userUC *userUsecase) UpdateUser(uid uint64, u *model.User) (*model.User, error) {
	updatedUser, err := userUC.userRepo.UpdateUser(uint64(uid), u)
	if err != nil {
		return &model.User{}, nil
	}

	return updatedUser, nil
}

func (userUC *userUsecase) DeleteUser(uid uint64) (int64, error) {
	deletedUserID, err := userUC.userRepo.DeleteUser(uint64(uid))
	if err != nil {
		return 0, err
	}
	return deletedUserID, nil
}

func (userUC *userUsecase) GetUserByID(uid uint64) (*model.User, error) {
	user, err := userUC.userRepo.GetUserByID(uint64(uid))
	if err != nil {
		return &model.User{}, nil
	}

	return user, nil
}

func (userUC *userUsecase) GetAllUsers() (*[]model.User, error) {
	users, err := userUC.userRepo.GetAllUsers()
	if err != nil {
		return &[]model.User{}, nil
	}

	return users, nil
}
