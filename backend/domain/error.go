package domain

type ErrorResponse struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}
