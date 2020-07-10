package encode

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sammy9867/daily-diary/backend/domain"
)

// JSON encodes the response in JSON format
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

// ERROR encodes the error response in JSON format
func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, domain.ErrorResponse{
			StatusCode: statusCode,
			Message:    err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}
