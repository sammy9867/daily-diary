package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
	"github.com/joho/godotenv"

	_userController "github.com/sammy9867/daily-diary/backend/user/controller"
	_userRepo "github.com/sammy9867/daily-diary/backend/user/repository/mysql"
	_userUseCase "github.com/sammy9867/daily-diary/backend/user/usecase/usecaseimpl"

	"github.com/sammy9867/daily-diary/backend/user/model"
)

func Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (DB *gorm.DB) {

	var err error

	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	}

	return DB

}

func run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	DB := Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	userRepo := _userRepo.NewMysqlUserRepository(DB)
	userUseCase := _userUseCase.NewUserUseCase(userRepo)

	DB.Debug().AutoMigrate(&model.User{}) //database migration

	router := mux.NewRouter()
	_userController.NewUserController(router, userUseCase)

	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", router))

}

func main() {

	run()
}
