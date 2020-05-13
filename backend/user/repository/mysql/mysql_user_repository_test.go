package mysql_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql driver
	"github.com/joho/godotenv"

	"github.com/go-playground/assert/v2"
	"github.com/sammy9867/daily-diary/backend/model"
	_userRepo "github.com/sammy9867/daily-diary/backend/user/repository/mysql"
)

type DatabaseConnection struct {
	DB *gorm.DB
}

var dbConn = DatabaseConnection{}

func TestCreateUser(t *testing.T) {

	mockDB()
	newUser := model.User{
		ID:       1,
		Username: "Sammy",
		Email:    "sammy@gmail.com",
		Password: "password",
	}

	savedUser, err := _userRepo.NewMysqlUserRepository(dbConn.DB).CreateUser(&newUser)
	if err != nil {
		t.Errorf("Error saving the user: %v\n", err)
		return
	}

	assert.Equal(t, newUser.ID, newUser.ID)
	assert.Equal(t, newUser.Username, savedUser.Username)
	assert.Equal(t, newUser.Email, savedUser.Email)

}

func TestUpdateUser(t *testing.T) {

	mockDB()
	user := model.User{
		ID:       1,
		Username: "Sammy",
		Email:    "sammy@gmail.com",
		Password: "password",
	}

	err := dbConn.DB.Model(&model.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Error saving the user: %v\n", err)
	}

	userUpdate := model.User{
		ID:       1,
		Username: "sammyUpdated",
		Email:    "sammyUpdated@gmail.com",
		Password: "password",
	}

	updatedUser, err := _userRepo.NewMysqlUserRepository(dbConn.DB).UpdateUser(userUpdate.ID, &userUpdate)
	if err != nil {
		t.Errorf("Error while updating the user: %v\n", err)
		return
	}

	assert.Equal(t, userUpdate.ID, updatedUser.ID)
	assert.Equal(t, userUpdate.Username, updatedUser.Username)
	assert.Equal(t, userUpdate.Email, updatedUser.Email)

}

func TestDeleteUser(t *testing.T) {

	mockDB()
	user := model.User{
		ID:       1,
		Username: "Sammy",
		Email:    "sammy@gmail.com",
		Password: "password",
	}

	err := dbConn.DB.Model(&model.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Error saving the user: %v\n", err)
	}

	isDeleted, err := _userRepo.NewMysqlUserRepository(dbConn.DB).DeleteUser(user.ID)
	if err != nil {
		t.Errorf("Error while deleting the user: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}

func TestGetUserByID(t *testing.T) {

	mockDB()
	user := model.User{
		ID:       1,
		Username: "Sammy",
		Email:    "sammy@gmail.com",
		Password: "password",
	}

	err := dbConn.DB.Model(&model.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Error saving the user: %v\n", err)
	}

	userFound, err := _userRepo.NewMysqlUserRepository(dbConn.DB).GetUserByID(user.ID)
	if err != nil {
		t.Errorf("Error while finding the user: %v\n", err)
		return
	}

	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Username, userFound.Username)
	assert.Equal(t, user.Email, userFound.Email)
}

func TestGetAllUsers(t *testing.T) {

	mockDB()
	users := []model.User{

		model.User{
			ID:       1,
			Username: "Sammy",
			Email:    "sammy@gmail.com",
			Password: "password",
		},
		model.User{
			ID:       2,
			Username: "pammy",
			Email:    "pammy@gmail.com",
			Password: "password",
		},
	}

	for i := range users {
		err := dbConn.DB.Model(&model.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("Error saving the user: %v\n", err)
		}
	}

	usersFound, err := _userRepo.NewMysqlUserRepository(dbConn.DB).GetAllUsers()
	if err != nil {
		t.Errorf("Error while finding users: %v\n", err)
		return
	}

	assert.Equal(t, len(users), len(*usersFound))
}

func mockDB() {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../../.env"))

	if err != nil {
		log.Fatalf("Error opening env %v\n", err)
	}

	dbConn.InitializeDBTest(os.Getenv("DB_DRIVER_TEST"), os.Getenv("DB_USER_TEST"), os.Getenv("DB_PASSWORD_TEST"), os.Getenv("DB_PORT_TEST"),
		os.Getenv("DB_HOST_TEST"), os.Getenv("DB_NAME_TEST"))

	err = dbConn.DB.DropTableIfExists(&model.User{}).Error
	if err != nil {
		log.Printf("Error dropping table %v\n", err)
	}

	err = dbConn.DB.AutoMigrate(&model.User{}).Error
	if err != nil {
		log.Printf("Error migrating User %v\n", err)
	}

	log.Printf("Successfully refreshed table")
}

func (dbConnec *DatabaseConnection) InitializeDBTest(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		dbConnec.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", Dbdriver)
		}
	}
}
