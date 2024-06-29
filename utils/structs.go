package utils

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Error struct {
	Code  int         `json:"code"`
	Error interface{} `json:"error"`
}
