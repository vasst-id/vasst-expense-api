package entities

type ApiResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta"`
}
