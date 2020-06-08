package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/sammy9867/daily-diary/backend/util/auth"
	"github.com/sammy9867/daily-diary/backend/util/encode"
)

var ApiKey = ""

func ApiKeySetter(apiKey string) {
	ApiKey = apiKey
}
func apiKeyGetter() string {
	return ApiKey
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

func init() {
	go cleanupVisitors()
}

func getVisitor(apiKey string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[apiKey]
	if !exists {
		limiter := rate.NewLimiter(1, 5) // Initial and maximum bucket size of 5 tokens, and 1 token is added every second.
		// Include the current time when creating a new visitor.
		visitors[apiKey] = &visitor{limiter, time.Now()}
		return limiter
	}

	// Update the last seen time for the visitor.
	v.lastSeen = time.Now()
	return v.limiter
}

// Every minute check the map for visitors that haven't been seen for
// more than 5 minutes and delete the entries.
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 5*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

// SetMiddlewareJSON will format all responses to JSON.
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		apiKey := apiKeyGetter()
		fmt.Println("ApiKey: ", apiKey)

		limiter := getVisitor(apiKey)
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

// SetMiddlewareAuthentication will check whether the user is authenticated or not
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.ValidateToken(r)
		if err != nil {
			encode.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
