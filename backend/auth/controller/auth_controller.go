package controller

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/sammy9867/daily-diary/backend/auth/controller/format"
	"github.com/sammy9867/daily-diary/backend/auth/usecase"
	"github.com/sammy9867/daily-diary/backend/domain"
	"github.com/sammy9867/daily-diary/backend/util/encode"
	"github.com/sammy9867/daily-diary/backend/util/middleware"
	"github.com/sammy9867/daily-diary/backend/util/token"
)

// AuthController represents all the http request of the users authentication
type AuthController struct {
	authUC usecase.AuthUseCase
	pool   *redis.Pool
}

// NewAuthController creates an object of AuthController
func NewAuthController(router *mux.Router, pool *redis.Pool, us usecase.AuthUseCase) {
	controller := &AuthController{
		authUC: us,
		pool:   pool,
	}
	router.HandleFunc("/login", middleware.SetMiddlewareJSON(controller.Login)).Methods("POST")
	router.HandleFunc("/logout", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(controller.Logout))).Methods("POST")
}

// Login endpoint
func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
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

	format.Initialize(&user)
	err = format.Validate(&user, "login")
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	tokenDetail, err := ac.authUC.Login(user.Email, user.Password)
	if err != nil {
		formattedError := format.FormatError(err.Error())
		encode.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	encode.JSON(w, http.StatusOK, tokenDetail)
}

// Logout endpoint
func (ac *AuthController) Logout(w http.ResponseWriter, r *http.Request) {

	authDetails, err := token.ExtractTokenMetaData(r)
	if err != nil {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	deleted, err := ac.authUC.Logout(authDetails.AccessUUID)
	if err != nil || deleted == 0 {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	encode.JSON(w, http.StatusOK, "Logged out")
}
