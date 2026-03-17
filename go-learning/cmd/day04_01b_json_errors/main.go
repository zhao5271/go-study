package main

import (
	"log"
	"net/http"
	"os"

	"example.com/go-learning/internal/httpkit"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type ListUsersData struct {
	Items []User `json:"items"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Total int    `json:"total"`
}

func main() {
	users := []User{
		{ID: 1, Name: "Alice", Role: "admin"},
		{ID: 2, Name: "Bob", Role: "viewer"},
		{ID: 3, Name: "Carol", Role: "editor"},
		{ID: 4, Name: "Dave", Role: "viewer"},
	}

	mux := http.NewServeMux()

	// 真实项目里建议统一加 /api/v1 做版本前缀，避免未来破坏性变更无处安放。
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httpkit.WriteError(w, http.StatusMethodNotAllowed, 10001, "METHOD_NOT_ALLOWED") // Output: {"code":10001,"message":"METHOD_NOT_ALLOWED"}\n
			return
		}
		httpkit.WriteJSON(w, http.StatusOK, httpkit.APIResponse{Code: 0, Message: "OK", Data: map[string]bool{"ok": true}}) // Output: {"code":0,"message":"OK","data":{"ok":true}}\n
	})

	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httpkit.WriteError(w, http.StatusMethodNotAllowed, 10001, "METHOD_NOT_ALLOWED") // Output: {"code":10001,"message":"METHOD_NOT_ALLOWED"}\n
			return
		}

		page, size, err := httpkit.ParsePageSize(r)
		if err != nil {
			httpkit.WriteError(w, http.StatusBadRequest, 10002, "INVALID_QUERY") // Output: {"code":10002,"message":"INVALID_QUERY"}\n
			return
		}

		total := len(users)
		start := (page - 1) * size
		if start > total {
			start = total
		}
		end := start + size
		if end > total {
			end = total
		}

		log.Printf("list_users page=%d size=%d start=%d end=%d", page, size, start, end) // Output: 2006/01/02 15:04:05 list_users page=1 size=2 start=0 end=2 (输出可能变化/不固定：包含时间戳)
		data := ListUsersData{
			Items: users[start:end],
			Page:  page,
			Size:  size,
			Total: total,
		}
		httpkit.WriteJSON(w, http.StatusOK, httpkit.APIResponse{Code: 0, Message: "OK", Data: data}) // Output: {"code":0,"message":"OK","data":{"items":[...],"page":1,"size":2,"total":4}}\n
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("listening on %s", addr) // Output: 2006/01/02 15:04:05 listening on :8080 (输出可能变化/不固定：包含时间戳)
	err := http.ListenAndServe(addr, mux)
	log.Printf("server stopped: %v", err) // Output: 2006/01/02 15:04:05 server stopped: listen tcp :8080: bind: address already in use (输出可能变化/不固定：取决于环境与错误)
}
