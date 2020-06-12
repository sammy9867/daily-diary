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
func ReJSONSet(rh *rejson.Handler, eid uint64, entry *domain.Entry) {

	fmt.Println("Redis SET")

	_, err := rh.JSONSet("entry:"+strconv.FormatUint(eid, 10), ".", entry)
	if err != nil {
		fmt.Println(err)
	}
}

// ReJSONGet gets the json from redis
func ReJSONGet(rh *rejson.Handler, uid uint64) (*domain.Entry, error) {

	fmt.Println("Redis GET")

	entryJSON, err := redis.Bytes(rh.JSONGet("entry:"+strconv.FormatUint(uid, 10), "."))
	if err != nil {
		return &domain.Entry{}, err
	}
	// Save returned value in entry struct
	var entry domain.Entry
	err = json.Unmarshal(entryJSON, &entry)
	if err != nil {
		return &domain.Entry{}, err
	}

	return &entry, nil
}

// ReJSONDel deletes the json from redis
func ReJSONDel(rh *rejson.Handler, eid uint64) {

	fmt.Println("Redis DEL")

	_, err := rh.JSONDel("entry:"+strconv.FormatUint(eid, 10), ".")
	if err != nil {
		fmt.Println(err)
	}
}
