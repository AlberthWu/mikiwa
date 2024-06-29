package utils

type ApiError struct {
	Param   string `json:"field"`
	Message string `json:"errors"`
}

type ResultChan struct {
	Id      int    `json:"id"`
	Data    string `json:"data"`
	Message string `json:"message"`
}
