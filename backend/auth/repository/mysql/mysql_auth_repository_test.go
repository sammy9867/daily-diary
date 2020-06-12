package mysql_test

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql driver
	"github.com/joho/godotenv"

	"github.com/go-playground/assert/v2"
	_authRepo "github.com/sammy9867/daily-diary/backend/auth/repository/mysql"
	"github.com/sammy9867/daily-diary/backend/domain"
)

type mockDatabaseConnection struct {
	DB   *gorm.DB
	pool *redis.Pool
}

var dbConn = mockDatabaseConnection{}

func TestLogin(t *testing.T) {

	mockDB()
	var err error
	user := domain.User{
		ID:       1,
		Username: "Sammy",
		Email:    "sammy@gmail.com",
		Password: "password",
	}

	err = dbConn.DB.Model(&domain.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Error saving the user: %v\n", err)
	}

	tokenDetails, err := _authRepo.NewMysqlAuthRepository(dbConn.DB, dbConn.pool).Login(user.Email, user.Password)
	if err != nil {
		t.Errorf("Error while login: %v\n", err)
		return
	}

	assert.NotEqual(t, tokenDetails.AccessToken, nil)
	assert.NotEqual(t, tokenDetails.RefreshToken, nil)
}

func TestLogout(t *testing.T) {

	mockDB()
	var err error
	user := domain.User{
		ID:       1,
		Username: "Sammy",
		Email:    "sammy@gmail.com",
		Password: "password",
	}

	err = dbConn.DB.Model(&domain.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Error saving the user: %v\n", err)
	}

	deleted, err := _authRepo.NewMysqlAuthRepository(dbConn.DB, dbConn.pool).Logout(strconv.FormatUint(user.ID, 10))
	if err != nil {
		t.Errorf("Error while logout: %v\n", err)
		return
	}

	assert.NotEqual(t, deleted, 1)
}

func TestRefresh(t *testing.T) {

	mockDB()
	var err error
	user := domain.User{
		ID:       1,
		Username: "Sammy",
		Email:    "sammy@gmail.com",
		Password: "password",
	}

	err = dbConn.DB.Model(&domain.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Error saving the user: %v\n", err)
	}

	// User login
	oldTokenDetails, err := _authRepo.NewMysqlAuthRepository(dbConn.DB, dbConn.pool).Login(user.Email, user.Password)
	if err != nil {
		t.Errorf("Error while login: %v\n", err)
		return
	}

	// Get refresh token from login
	refreshToken := oldTokenDetails.RefreshToken
	newTokenDetails, err := _authRepo.NewMysqlAuthRepository(dbConn.DB, dbConn.pool).Refresh(refreshToken)
	if err != nil {
		t.Errorf("Error while refreshing the token: %v\n", err)
		return
	}

	assert.NotEqual(t, newTokenDetails.AccessToken, nil)
	assert.NotEqual(t, newTokenDetails.RefreshToken, nil)
}

func mockDB() {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../../.env"))

	if err != nil {
		log.Fatalf("Error opening env %v\n", err)
	}

	dbConn.InitializeDBTest(os.Getenv("DB_DRIVER_TEST"), os.Getenv("DB_USER_TEST"), os.Getenv("DB_PASSWORD_TEST"), os.Getenv("DB_PORT_TEST"),
		os.Getenv("DB_HOST_TEST"), os.Getenv("DB_NAME_TEST"))

	dbConn.InitializeRedisCacheTest(10, "localhost:6379")

	if err := dbConn.DB.Raw("CALL TrucateTables()").Scan(&domain.EntryImage{}).Scan(&domain.Entry{}).Scan(&domain.User{}).Error; err != nil {
		log.Printf("Error truncating tables: %v\n", err)
	}

	log.Printf("Successfully refreshed table")
}

func (dbConnec *mockDatabaseConnection) InitializeDBTest(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

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

func (dbConnec *mockDatabaseConnection) InitializeRedisCacheTest(maxIdleConn int, port string) {

	dbConnec.pool = &redis.Pool{
		MaxIdle:     maxIdleConn,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", port)
		},
	}

	conn := dbConnec.pool.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHALL")
	if err != nil {
		fmt.Printf("Could not flush data from redis server")
	}
}
