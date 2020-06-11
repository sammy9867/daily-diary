package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
	"github.com/joho/godotenv"
	"github.com/nitishm/go-rejson"

	_userController "github.com/sammy9867/daily-diary/backend/user/controller"
	_userRepo "github.com/sammy9867/daily-diary/backend/user/repository/mysql"
	_userUseCase "github.com/sammy9867/daily-diary/backend/user/usecase/usecaseimpl"

	_entryController "github.com/sammy9867/daily-diary/backend/entry/controller"
	_entryRepo "github.com/sammy9867/daily-diary/backend/entry/repository/mysql"
	_entryUseCase "github.com/sammy9867/daily-diary/backend/entry/usecase/usecaseimpl"
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
			fmt.Printf("We are connected to the %s database\n", Dbdriver)
		}
	}

	return DB

}

func InitializeRedisCache(maxIdleConn int, port string) *redis.Pool {

	pool := &redis.Pool{
		MaxIdle:     maxIdleConn,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", port)
		},
	}

	return pool
}

func run() {

	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error opening env, %v", err)
	} else {
		fmt.Println(".env file loaded")
	}

	DB := Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	cachePool := InitializeRedisCache(10, "localhost:6379")
	conn := cachePool.Get()
	defer conn.Close()
	rh := rejson.NewReJSONHandler()
	rh.SetRedigoClient(conn)

	userRepo := _userRepo.NewMysqlUserRepository(DB, rh)
	userUseCase := _userUseCase.NewUserUseCase(userRepo)

	entryRepo := _entryRepo.NewMysqlEntryRepository(DB)
	entryUseCase := _entryUseCase.NewEntryUseCase(entryRepo)

	router := mux.NewRouter()
	_userController.NewUserController(router, userUseCase)
	_entryController.NewEntryController(router, entryUseCase)

	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", router))

}

func main() {

	run()
}
