package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/sammy9867/daily-diary/backend/domain"
	"github.com/sammy9867/daily-diary/backend/user/controller/format"
	"github.com/sammy9867/daily-diary/backend/user/usecase"
	"github.com/sammy9867/daily-diary/backend/util/encode"
	"github.com/sammy9867/daily-diary/backend/util/middleware"
	"github.com/sammy9867/daily-diary/backend/util/token"
)

// UserController represents all the http request of the user
type UserController struct {
	userUC usecase.UserUseCase
	pool   *redis.Pool
}

// NewUserController creates an object of UserController
func NewUserController(router *mux.Router, pool *redis.Pool, us usecase.UserUseCase) {
	controller := &UserController{
		userUC: us,
		pool:   pool,
	}

	router.HandleFunc("/users", middleware.SetMiddlewareJSON(controller.CreateUser)).Methods("POST")
	router.HandleFunc("/users/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(controller.UpdateUser))).Methods("PUT")
	router.HandleFunc("/users/{id}", middleware.SetMiddlewareAuthentication(controller.DeleteUser)).Methods("DELETE")
	router.HandleFunc("/users/{id}", middleware.SetMiddlewareJSON(controller.GetUserByID)).Methods("GET")
	router.HandleFunc("/users", middleware.SetMiddlewareJSON(controller.GetAllUsers)).Methods("GET")
}

// CreateUser endpoint is used to create a new user using usersname, email and password
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := domain.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Initialize()
	err = format.Validate(&user, "")
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	createdUser, err := uc.userUC.CreateUser(&user)
	if err != nil {
		checkError := format.CheckError(err.Error())
		encode.ERROR(w, http.StatusConflict, checkError)
		fmt.Println("CreateUser", checkError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, createdUser.ID))
	encode.JSON(w, http.StatusCreated, createdUser)
}

// UpdateUser endpoint is used to update a user's credentials
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		encode.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := domain.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	authDetails, err := token.ExtractTokenMetaData(r)
	fmt.Println(authDetails)
	if err != nil {
		fmt.Println(1, err)
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	userID, err := token.FetchAuthDetails(authDetails, uc.pool)
	if err != nil {
		fmt.Println(2, err)
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if userID != uint64(uid) {
		fmt.Println(3, err)
		encode.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	user.Initialize()
	err = format.Validate(&user, "update")
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedUser, err := uc.userUC.UpdateUser(uint64(uid), &user)
	if err != nil {
		encode.ERROR(w, http.StatusInternalServerError, errors.New("Error while updating the user"))
		fmt.Println("UpdateUser", err)
		return
	}
	encode.JSON(w, http.StatusOK, updatedUser)
}

// DeleteUser endpoint
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		encode.ERROR(w, http.StatusBadRequest, err)
		return
	}

	authDetails, err := token.ExtractTokenMetaData(r)
	if err != nil {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	userID, err := token.FetchAuthDetails(authDetails, uc.pool)
	if err != nil {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if userID != 0 && userID != uint64(uid) {
		encode.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	_, err = uc.userUC.DeleteUser(uint64(uid))
	if err != nil {
		encode.ERROR(w, http.StatusInternalServerError, errors.New("Error while deleting the user"))
		fmt.Println("DeleteUser", err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	encode.JSON(w, http.StatusNoContent, "")
}

// GetUserByID endpoint
func (uc *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		encode.ERROR(w, http.StatusBadRequest, err)
		return
	}
	userGotten, err := uc.userUC.GetUserByID(uint64(uid))
	if err != nil {
		encode.ERROR(w, http.StatusInternalServerError, err)
		fmt.Println("GetUserByID", err)
		return
	}
	encode.JSON(w, http.StatusOK, userGotten)
}

// GetAllUsers endpoint
func (uc *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := uc.userUC.GetAllUsers()
	if err != nil {
		encode.ERROR(w, http.StatusInternalServerError, err)
		fmt.Println("GetAllUsers", err)
		return
	}
	encode.JSON(w, http.StatusOK, users)
}
