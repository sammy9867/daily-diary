package mysql

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/sammy9867/daily-diary/backend/auth/repository"
	"github.com/sammy9867/daily-diary/backend/domain"
	"github.com/sammy9867/daily-diary/backend/util/token"
	"golang.org/x/crypto/bcrypt"
)

type mysqlAuthRepository struct {
	DB   *gorm.DB
	pool *redis.Pool
}

// NewMysqlAuthRepository will create an object that will implement AuthRepository interface
// Note: Need to implement all the methods from the interface
func NewMysqlAuthRepository(DB *gorm.DB, pool *redis.Pool) repository.AuthRepository {
	return &mysqlAuthRepository{DB, pool}
}

func (mysqlAuthRepo *mysqlAuthRepository) Login(email, password string) (*domain.TokenDetail, error) {

	var err error

	user := domain.User{}

	err = mysqlAuthRepo.DB.Debug().Model(domain.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}
	err = VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, err
	}

	tokenDetails, err := token.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}

	err = token.SaveTokenMetaData(user.ID, tokenDetails, mysqlAuthRepo.pool)
	if err != nil {
		return nil, err
	}

	return tokenDetails, nil
}

// func (mysqlAuthRepo *mysqlAuthRepository) Logout() error {
// 	return nil
// }

// VerifyPassword will check if the user's password matched with the hashed password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
