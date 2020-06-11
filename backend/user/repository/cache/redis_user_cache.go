package cache

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson"
	"github.com/sammy9867/daily-diary/backend/domain"
)

// ReJsonSet saves the json in redis
func ReJsonSet(rh *rejson.Handler, uid uint64, user *domain.User) {

	fmt.Println("Redis SET")

	_, err := rh.JSONSet("user:"+strconv.FormatUint(uid, 10), ".", user)
	if err != nil {
		fmt.Println(err)
	}
}

// ReJsonGet gets the json from redis
func ReJsonGet(rh *rejson.Handler, uid uint64) (*domain.User, error) {

	fmt.Println("Redis GET")

	userJSON, err := redis.Bytes(rh.JSONGet("user:"+strconv.FormatUint(uid, 10), "."))
	if err != nil {
		return &domain.User{}, err
	}
	// Save returned value in user struct
	var user domain.User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return &domain.User{}, err
	}

	return &user, nil
}

// ReJsonDel deletes the json from redis
func ReJsonDel(rh *rejson.Handler, uid uint64) {

	fmt.Println("Redis DEL")

	_, err := rh.JSONDel("user:"+strconv.FormatUint(uid, 10), ".")
	if err != nil {
		fmt.Println(err)
	}
}
