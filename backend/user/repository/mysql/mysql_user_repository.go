package mysql

import (
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sammy9867/daily-diary/backend/user/controller/auth"
	"github.com/sammy9867/daily-diary/backend/user/model"
	"github.com/sammy9867/daily-diary/backend/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type mysqlUserRepository struct {
	DB *gorm.DB
}

// NewMysqlUserRepository will create an object that will implement UserRepository interface
// Note: Need to implement all the methods from the interface
func NewMysqlUserRepository(DB *gorm.DB) repository.UserRepository {
	return &mysqlUserRepository{DB}
}

func (mysqlUserRepo *mysqlUserRepository) SignIn(email, password string) (string, error) {

	var err error

	user := model.User{}

	err = mysqlUserRepo.DB.Debug().Model(model.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}

func (mysqlUserRepo *mysqlUserRepository) CreateUser(u *model.User) (*model.User, error) {

	var err error
	err = BeforeSave(u)
	if err != nil {
		log.Fatal(err)
	}
	err = mysqlUserRepo.DB.Debug().Create(&u).Error
	if err != nil {
		return &model.User{}, err
	}

	return u, nil
}

func (mysqlUserRepo *mysqlUserRepository) UpdateUser(uid uint64, u *model.User) (*model.User, error) {

	err := BeforeSave(u)
	if err != nil {
		log.Fatal(err)
	}

	db := mysqlUserRepo.DB.Debug().Model(&model.User{}).Where("id = ?", uid).UpdateColumns(
		map[string]interface{}{
			"username":   u.Username,
			"email":      u.Email,
			"password":   u.Password,
			"updated_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &model.User{}, db.Error
	}

	user, err := mysqlUserRepo.GetUserByID(uid)
	if err != nil {
		return &model.User{}, err
	}
	return user, nil
}

func (mysqlUserRepo *mysqlUserRepository) DeleteUser(uid uint64) (int64, error) {
	db := mysqlUserRepo.DB.Debug().Model(&model.User{}).Where("id = ?", uid).Take(&model.User{}).Delete(&model.User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (mysqlUserRepo *mysqlUserRepository) GetUserByID(uid uint64) (*model.User, error) {

	var err error
	user := model.User{}
	err = mysqlUserRepo.DB.Debug().Model(model.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		return &model.User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &model.User{}, errors.New("User Not Found")
	}
	return &user, err
}

func (mysqlUserRepo *mysqlUserRepository) GetAllUsers() (*[]model.User, error) {

	var err error
	users := []model.User{}

	err = mysqlUserRepo.DB.Debug().Model(&model.User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]model.User{}, err
	}
	return &users, err

}

// Hash will hash the user's password
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword will check if the user's password matched with the hashed password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// BeforeSave will hash the password before creating/updating a user
func BeforeSave(u *model.User) error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
