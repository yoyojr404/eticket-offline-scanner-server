package ws

import "encoding/json"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func Success(data []byte, dest *Response) {
	dest.Success = true
	dest.Message = "OK"
	if len(data) > 0 {
		dest.Data = json.RawMessage(data)
	}
}

func Error(err error, dest *Response) {
	dest.Success = false
	dest.Message = "FAILED"
	dest.Error = err.Error()
}
