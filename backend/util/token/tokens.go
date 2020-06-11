package token

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sammy9867/daily-diary/backend/domain"
	uuid "github.com/satori/go.uuid"

	jwt "github.com/dgrijalva/jwt-go"
)

// CreateToken will create an access token valid for 5 minutues and a refresh token valid for 7 days using HS256 algorithm
func CreateToken(userID uint64) (*domain.TokenDetail, error) {
	td := &domain.TokenDetail{}
	td.AcctokenExpiresAt = time.Now().Add(time.Minute * 5).Unix()
	td.AccessUUID = uuid.NewV4().String()

	td.ReftokenExpiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.NewV4().String()

	var err error
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["access_uuid"] = td.AccessUUID
	accessTokenClaims["user_id"] = userID
	accessTokenClaims["exp"] = td.AcctokenExpiresAt
	actoken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	td.AccessToken, err = actoken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["refresh_uuid"] = td.RefreshUUID
	refreshTokenClaims["user_id"] = userID
	refreshTokenClaims["exp"] = td.ReftokenExpiresAt
	retoken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	td.RefreshToken, err = retoken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

// SaveTokenMetaData saves the access and refresh token in redis
func SaveTokenMetaData(userID uint64, td *domain.TokenDetail, pool *redis.Pool) error {
	at := time.Unix(td.AcctokenExpiresAt, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.ReftokenExpiresAt, 0)

	conn := pool.Get()
	defer conn.Close()

	fmt.Println(td.AccessUUID, strconv.Itoa(int(userID)), math.Ceil(at.Sub(time.Now()).Seconds()))
	fmt.Println(td.RefreshUUID, strconv.Itoa(int(userID)), math.Ceil(rt.Sub(time.Now()).Seconds()))

	var err error
	_, err = conn.Do("SET", td.AccessUUID, strconv.Itoa(int(userID)), "EX", math.Ceil(at.Sub(time.Now()).Seconds()))
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", td.RefreshUUID, strconv.Itoa(int(userID)), "EX", math.Ceil(rt.Sub(time.Now()).Seconds()))
	if err != nil {
		return err
	}

	return nil
}

// ValidateToken will validate the token passed in the authentication header. Called in middleware/middleware.go
func ValidateToken(r *http.Request) error {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return err
	}
	return nil
}

// ExtractToken from the request header
func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// ExtractTokenMetaData will extract the authenication details that contain accessUUID and userID
func ExtractTokenMetaData(r *http.Request) (*domain.AuthDetail, error) {

	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &domain.AuthDetail{
			AccessUUID: accessUUID,
			UserID:     uid,
		}, nil
	}
	return nil, err
}

// FetchAuthDetails fetches authentication details
func FetchAuthDetails(authDetail *domain.AuthDetail, pool *redis.Pool) (uint64, error) {

	conn := pool.Get()
	defer conn.Close()

	userid, err := redis.String(conn.Do("GET", authDetail.AccessUUID))
	if err != nil {
		return 0, err
	}

	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

// DeleteAuth deletes uuid from redis for logout
func DeleteAuth(uuid string, pool *redis.Pool) (int64, error) {
	conn := pool.Get()
	defer conn.Close()

	deleted, err := redis.Int64(conn.Do("DEL", uuid))
	if err != nil {
		return 0, err
	}

	return deleted, nil
}
