package utils

type ApiError struct {
	Param   string `json:"field"`
	Message string `json:"errors"`
}
