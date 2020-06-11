package cache

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson"
	"github.com/sammy9867/daily-diary/backend/domain"
)

// ReJSONSet saves the json in redis
func ReJSONSet(rh *rejson.Handler, uid uint64, user *domain.User) {

	fmt.Println("Redis SET")

	_, err := rh.JSONSet("user:"+strconv.FormatUint(uid, 10), ".", user)
	if err != nil {
		fmt.Println(err)
	}
}

// ReJSONGet gets the json from redis
func ReJSONGet(rh *rejson.Handler, uid uint64) (*domain.User, error) {

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

// ReJSONDel deletes the json from redis
func ReJSONDel(rh *rejson.Handler, uid uint64) {

	fmt.Println("Redis DEL")

	_, err := rh.JSONDel("user:"+strconv.FormatUint(uid, 10), ".")
	if err != nil {
		fmt.Println(err)
	}
}
