package httpkit

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v) // Output: 输出可能变化/不固定（取决于 v 的内容）
}

func WriteError(w http.ResponseWriter, status int, code int, message string) {
	WriteJSON(w, status, APIResponse{Code: code, Message: message})
}

