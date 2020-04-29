package mysql

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/sammy9867/daily-diary/backend/user/model"
	"github.com/sammy9867/daily-diary/backend/user/repository"
)

type mysqlUserRepository struct {
	DB *gorm.DB
}

// NewMysqlUserRepository will create an object that will implement UserRepository interface
// Need to implement all the methods from the interface
func NewMysqlUserRepository(DB *gorm.DB) repository.UserRepository {
	return &mysqlUserRepository{DB}
}

func (mysqlUserRepo *mysqlUserRepository) CreateUser(u *model.User) (*model.User, error) {

	var err error
	err = mysqlUserRepo.DB.Debug().Create(&u).Error
	if err != nil {
		return &model.User{}, nil
	}

	return u, nil
}

func (mysqlUserRepo *mysqlUserRepository) UpdateUser(uid uint64, u *model.User) (*model.User, error) {

	db := mysqlUserRepo.DB.Debug().Model(&model.User{}).Where("id = ?", uid).Take(&model.User{}).UpdateColumns(
		map[string]interface{}{
			"password":   u.Password,
			"username":   u.Username,
			"email":      u.Email,
			"updated_at": u.UpdatedAt,
		},
	)
	if db.Error != nil {
		return &model.User{}, db.Error
	}

	// This is the display the updated user
	err := mysqlUserRepo.DB.Debug().Model(&model.User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &model.User{}, err
	}
	return u, nil
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
