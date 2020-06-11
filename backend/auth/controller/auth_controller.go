package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sammy9867/daily-diary/backend/auth/controller/format"
	"github.com/sammy9867/daily-diary/backend/auth/usecase"
	"github.com/sammy9867/daily-diary/backend/domain"
	"github.com/sammy9867/daily-diary/backend/util/encode"
	"github.com/sammy9867/daily-diary/backend/util/middleware"
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
	router.HandleFunc("/login", middleware.SetMiddlewareJSON(controller.Login)).Methods("POST")
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
