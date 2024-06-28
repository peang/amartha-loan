package utils

type Response struct {
	Data interface{} `json:"data,omitempty"`
}

type Error struct {
	Code  int         `json:"code"`
	Error interface{} `json:"error"`
}
