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

// Exercise 2 Answer:
// - 统一 404 JSON：用 ServeMux.Handler 做一次“是否命中路由”的探测
func withJSONNotFound(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, pattern := mux.Handler(r)
		if pattern == "" || h == nil {
			writeError(w, http.StatusNotFound, 10004, "NOT_FOUND") // Output: {"code":10004,"message":"NOT_FOUND"}\n
			return
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, APIResponse{Code: 0, Message: "OK", Data: map[string]bool{"ok": true}}) // Output: {"code":0,"message":"OK","data":{"ok":true}}\n
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("listening on %s", addr) // Output: 2006/01/02 15:04:05 listening on :8080 (输出可能变化/不固定：包含时间戳)
	err := http.ListenAndServe(addr, withJSONNotFound(mux))
	log.Printf("server stopped: %v", err) // Output: 2006/01/02 15:04:05 server stopped: listen tcp :8080: bind: address already in use (输出可能变化/不固定：取决于环境与错误)
}
