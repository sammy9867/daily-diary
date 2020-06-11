package mysql

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nitishm/go-rejson"
	"github.com/sammy9867/daily-diary/backend/domain"
	"github.com/sammy9867/daily-diary/backend/user/repository"
	"github.com/sammy9867/daily-diary/backend/user/repository/cache"
	"golang.org/x/crypto/bcrypt"
)

type mysqlUserRepository struct {
	DB *gorm.DB
	rh *rejson.Handler
}

// NewMysqlUserRepository will create an object that will implement UserRepository interface
// Note: Need to implement all the methods from the interface
func NewMysqlUserRepository(DB *gorm.DB, rh *rejson.Handler) repository.UserRepository {
	return &mysqlUserRepository{DB, rh}
}

func (mysqlUserRepo *mysqlUserRepository) CreateUser(u *domain.User) (*domain.User, error) {

	var err error
	err = BeforeSave(u)
	if err != nil {
		log.Fatal(err)
	}
	err = mysqlUserRepo.DB.Debug().Create(&u).Error
	if err != nil {
		return &domain.User{}, err
	}

	return u, nil
}

func (mysqlUserRepo *mysqlUserRepository) UpdateUser(uid uint64, u *domain.User) (*domain.User, error) {

	err := BeforeSave(u)
	if err != nil {
		log.Fatal(err)
	}

	db := mysqlUserRepo.DB.Debug().Model(&domain.User{}).Where("id = ?", uid).UpdateColumns(
		map[string]interface{}{
			"username":   u.Username,
			"email":      u.Email,
			"password":   u.Password,
			"updated_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &domain.User{}, db.Error
	}

	// Update Cache
	userCache, err := cache.ReJSONGet(mysqlUserRepo.rh, uid)
	if err != nil {
		return &domain.User{}, err
	}

	userCache.Username = u.Username
	userCache.Email = u.Email
	userCache.Password = u.Password
	userCache.UpdatedAt = time.Now()
	cache.ReJSONSet(mysqlUserRepo.rh, uid, userCache)

	user, err := mysqlUserRepo.GetUserByID(uid)
	if err != nil {
		return &domain.User{}, err
	}
	return user, nil
}

func (mysqlUserRepo *mysqlUserRepository) DeleteUser(uid uint64) (int64, error) {
	db := mysqlUserRepo.DB.Debug().Model(&domain.User{}).Where("id = ?", uid).Take(&domain.User{}).Delete(&domain.User{})

	if db.Error != nil {
		return 0, db.Error
	}

	// Delete Cache
	_, err := cache.ReJSONGet(mysqlUserRepo.rh, uid)
	if err != nil {
		return 0, err
	}
	cache.ReJSONDel(mysqlUserRepo.rh, uid)

	return db.RowsAffected, nil
}

func (mysqlUserRepo *mysqlUserRepository) GetUserByID(uid uint64) (*domain.User, error) {

	var err error

	// Check whether data exists in Redis
	user, err := cache.ReJSONGet(mysqlUserRepo.rh, uid)
	if err != nil {
		fmt.Println("fetch from db")
		err = mysqlUserRepo.DB.Debug().Model(domain.User{}).Where("id = ?", uid).Take(&user).Error
		if err != nil {
			return &domain.User{}, err
		}
		if gorm.IsRecordNotFoundError(err) {
			return &domain.User{}, errors.New("User Not Found")
		}

		// Update Redis
		cache.ReJSONSet(mysqlUserRepo.rh, uid, user)
	}

	return user, err
}

func (mysqlUserRepo *mysqlUserRepository) GetAllUsers() (*[]domain.User, error) {

	var err error
	users := []domain.User{}

	err = mysqlUserRepo.DB.Debug().Model(&domain.User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]domain.User{}, err
	}
	return &users, err

}

// Hash will hash the user's password
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// BeforeSave will hash the password before creating/updating a user
func BeforeSave(u *domain.User) error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
