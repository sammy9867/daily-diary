package usecaseimpl

import (
	"github.com/sammy9867/daily-diary/backend/auth/repository"
	"github.com/sammy9867/daily-diary/backend/auth/usecase"
	"github.com/sammy9867/daily-diary/backend/domain"
)

type authUsecase struct {
	authRepo repository.AuthRepository
}

// NewAuthUseCase will create an object that will implement AuthUseCase interface
// Note: Need to implement all the methods from the interface
func NewAuthUseCase(ar repository.AuthRepository) usecase.AuthUseCase {
	return &authUsecase{authRepo: ar}
}

func (authUC *authUsecase) Login(email, password string) (*domain.TokenDetail, error) {
	tokenDetail, err := authUC.authRepo.Login(email, password)
	return tokenDetail, err
}

func (authUC *authUsecase) Logout(uuid string) (int64, error) {
	deleted, err := authUC.authRepo.Logout(uuid)
	return deleted, err
}
