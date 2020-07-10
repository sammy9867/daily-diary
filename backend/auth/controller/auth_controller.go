package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

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
}

// NewAuthController creates an object of AuthController
func NewAuthController(router *mux.Router, us usecase.AuthUseCase) {
	controller := &AuthController{
		authUC: us,
	}
	router.HandleFunc("/auth/login", middleware.SetMiddlewareJSON(controller.Login)).Methods("POST")
	router.HandleFunc("/auth/logout", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(controller.Logout))).Methods("POST")
	router.HandleFunc("/auth/token/refresh", middleware.SetMiddlewareJSON(controller.RefreshToken)).Methods("POST")
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

	user.Initialize()
	err = format.Validate(&user, "login")
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	tokenDetail, err := ac.authUC.Login(user.Email, user.Password)
	if err != nil {
		encode.ERROR(w, http.StatusUnauthorized, errors.New("Incorrect Email or Password"))
		fmt.Println("Login", err)
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
		encode.ERROR(w, http.StatusInternalServerError, errors.New("Error while logging out"))
		fmt.Println("Logout", err)
		return
	}

	encode.JSON(w, http.StatusOK, "Logged out")
}

// RefreshToken endpoint
func (ac *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	refreshTokenMap := map[string]string{}
	err = json.Unmarshal(body, &refreshTokenMap)
	if err != nil {
		encode.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	refreshToken := refreshTokenMap["refresh_token"]
	tokenDetail, err := ac.authUC.Refresh(refreshToken)
	if err != nil {
		encode.ERROR(w, http.StatusInternalServerError, errors.New("Error while generating refresh token"))
		fmt.Println("RefreshToken", err)
		return
	}

	encode.JSON(w, http.StatusOK, tokenDetail)
}
