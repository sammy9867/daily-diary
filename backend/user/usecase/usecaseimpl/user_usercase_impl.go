package usecaseimpl

import (
	"github.com/sammy9867/daily-diary/backend/domain"
	"github.com/sammy9867/daily-diary/backend/user/repository"
	"github.com/sammy9867/daily-diary/backend/user/usecase"
)

type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserUseCase will create an object that will implement UserUserCase interface
// Note: Need to implement all the methods from the interface
func NewUserUseCase(ur repository.UserRepository) usecase.UserUseCase {
	return &userUsecase{userRepo: ur}
}

func (userUC *userUsecase) CreateUser(u *domain.User) (*domain.User, error) {
	createdUser, err := userUC.userRepo.CreateUser(u)
	return createdUser, err
}

func (userUC *userUsecase) UpdateUser(uid uint64, u *domain.User) (*domain.User, error) {
	updatedUser, err := userUC.userRepo.UpdateUser(uint64(uid), u)
	return updatedUser, err
}

func (userUC *userUsecase) DeleteUser(uid uint64) (int64, error) {
	deletedUserID, err := userUC.userRepo.DeleteUser(uint64(uid))
	return deletedUserID, err
}

func (userUC *userUsecase) GetUserByID(uid uint64) (*domain.User, error) {
	user, err := userUC.userRepo.GetUserByID(uint64(uid))
	return user, err
}

func (userUC *userUsecase) GetAllUsers() (*[]domain.User, error) {
	users, err := userUC.userRepo.GetAllUsers()
	return users, err
}
