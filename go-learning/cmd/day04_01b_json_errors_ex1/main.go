package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v) // Output: {"code":0,"message":"OK","data":{...}}\n
}

func writeError(w http.ResponseWriter, status int, code int, message string) {
	writeJSON(w, status, APIResponse{Code: code, Message: message})
}

// Exercise 1 Answer:
// - 新增中间件式的 method guard：只允许 GET，否则返回 405 JSON 错误
func onlyGET(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, 10001, "METHOD_NOT_ALLOWED") // Output: {"code":10001,"message":"METHOD_NOT_ALLOWED"}\n
			return
		}
		next(w, r)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/ping", onlyGET(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, APIResponse{Code: 0, Message: "OK", Data: map[string]string{"pong": "pong"}}) // Output: {"code":0,"message":"OK","data":{"pong":"pong"}}\n
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("listening on %s", addr) // Output: 2006/01/02 15:04:05 listening on :8080 (输出可能变化/不固定：包含时间戳)
	err := http.ListenAndServe(addr, mux)
	log.Printf("server stopped: %v", err) // Output: 2006/01/02 15:04:05 server stopped: listen tcp :8080: bind: address already in use (输出可能变化/不固定：取决于环境与错误)
}
