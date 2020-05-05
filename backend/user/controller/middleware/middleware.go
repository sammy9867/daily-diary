package middleware

import (
	"errors"
	"net/http"

	"github.com/sammy9867/daily-diary/backend/user/controller/auth"
	"github.com/sammy9867/daily-diary/backend/user/controller/util"
)

// SetMiddlewareJSON will format all responses to JSON.
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

//SetMiddlewareAuthentication will check whether the user is authenticated or not
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.ValidateToken(r)
		if err != nil {
			util.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
