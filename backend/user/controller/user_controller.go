package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sammy9867/daily-diary/backend/user/controller/auth"
	"github.com/sammy9867/daily-diary/backend/user/controller/middleware"
	"github.com/sammy9867/daily-diary/backend/user/controller/util"
	"github.com/sammy9867/daily-diary/backend/user/model"
	"github.com/sammy9867/daily-diary/backend/user/usecase"
)

type UserController struct {
	userUC usecase.UserUseCase
}

// NewUserController creates an object of UserController
func NewUserController(router *mux.Router, us usecase.UserUseCase) {
	controller := &UserController{
		userUC: us,
	}

	router.HandleFunc("/login", middleware.SetMiddlewareJSON(controller.Login)).Methods("POST")

	router.HandleFunc("/users", middleware.SetMiddlewareJSON(controller.CreateUser)).Methods("POST")
	router.HandleFunc("/users/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(controller.UpdateUser))).Methods("PUT") // User is not updated, returns id 0 for non-existing user
	router.HandleFunc("/users/{id}", middleware.SetMiddlewareAuthentication(controller.DeleteUser)).Methods("DELETE")                            //  returns id 0 after deletion
	router.HandleFunc("/users/{id}", middleware.SetMiddlewareJSON(controller.GetUserByID)).Methods("GET")                                        // Returns Empty User when id doesnt exist
	router.HandleFunc("/users", middleware.SetMiddlewareJSON(controller.GetAllUsers)).Methods("GET")                                             //  returns id 0 for non-existing user
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := model.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		util.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := uc.userUC.SignIn(user.Email, user.Password)
	if err != nil {
		util.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	util.JSON(w, http.StatusOK, token)
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := model.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		util.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	createdUser, err := uc.userUC.CreateUser(&user)
	if err != nil {
		util.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, createdUser.ID))
	util.JSON(w, http.StatusCreated, createdUser)
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := model.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		util.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		util.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint64(uid) {
		util.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	updatedUser, err := uc.userUC.UpdateUser(uint64(uid), &user)
	if err != nil {
		util.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	util.JSON(w, http.StatusOK, updatedUser)
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		util.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint64(uid) {
		util.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	_, err = uc.userUC.DeleteUser(uint64(uid))
	if err != nil {
		util.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	util.JSON(w, http.StatusNoContent, "")
}

func (uc *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, err)
		return
	}
	userGotten, err := uc.userUC.GetUserByID(uint64(uid))
	if err != nil {
		util.ERROR(w, http.StatusBadRequest, err)
		return
	}
	util.JSON(w, http.StatusOK, userGotten)
}

func (uc *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := uc.userUC.GetAllUsers()
	if err != nil {
		util.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	util.JSON(w, http.StatusOK, users)
}
