package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// 统一返回结构：前端更好做提示/重试/埋点；后端更好做告警与排障
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type ListUsersData struct {
	Items []User `json:"items"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Total int    `json:"total"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v) // Output: {"code":0,"message":"OK","data":{...}}\n (输出可能变化/不固定：取决于 v)
}

func writeError(w http.ResponseWriter, status int, code int, message string) {
	writeJSON(w, status, APIResponse{Code: code, Message: message}) // Output: {"code":10004,"message":"NOT_FOUND"}\n
}

func onlyMethod(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			writeError(w, http.StatusMethodNotAllowed, 10001, "METHOD_NOT_ALLOWED") // Output: {"code":10001,"message":"METHOD_NOT_ALLOWED"}\n
			return
		}
		next(w, r)
	}
}

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

var errInvalidQuery = errors.New("invalid query")

func parsePageSize(r *http.Request) (page int, size int, err error) {
	page = 1
	size = 20

	q := r.URL.Query()
	if raw := strings.TrimSpace(q.Get("page")); raw != "" {
		n, convErr := strconv.Atoi(raw)
		if convErr != nil || n < 1 {
			return 0, 0, errInvalidQuery
		}
		page = n
	}
	if raw := strings.TrimSpace(q.Get("size")); raw != "" {
		n, convErr := strconv.Atoi(raw)
		if convErr != nil || n < 1 || n > 100 {
			return 0, 0, errInvalidQuery
		}
		size = n
	}
	return page, size, nil
}

func main() {
	mux := http.NewServeMux()

	// 健康检查：给部署/容器编排做 readiness/liveness 用
	mux.HandleFunc("/api/v1/health", onlyMethod(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, APIResponse{Code: 0, Message: "OK", Data: map[string]bool{"ok": true}}) // Output: {"code":0,"message":"OK","data":{"ok":true}}\n
	}))

	// 用户列表：后台管理最常见的“列表分页检索”入口（先用假数据，后面再接 MySQL）
	mux.HandleFunc("/api/v1/users", onlyMethod(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		page, size, err := parsePageSize(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, 10002, "INVALID_QUERY") // Output: {"code":10002,"message":"INVALID_QUERY"}\n
			return
		}

		items := []User{
			{ID: 1, Email: "admin@example.com", Name: "Admin", Role: "admin"},
			{ID: 2, Email: "dev@example.com", Name: "Dev", Role: "dev"},
		}

		// 演示用：假装我们总共有 2 条（后面接 MySQL 时，这里会变成 COUNT(*)）
		total := 2
		log.Printf("list_users page=%d size=%d", page, size) // Output: 2006/01/02 15:04:05 list_users page=1 size=20 (输出可能变化/不固定：包含时间戳/参数变化)

		writeJSON(w, http.StatusOK, APIResponse{
			Code:    0,
			Message: "OK",
			Data: ListUsersData{
				Items: items,
				Page:  page,
				Size:  size,
				Total: total,
			},
		}) // Output: {"code":0,"message":"OK","data":{"items":[...],"page":1,"size":20,"total":2}}\n
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "18080"
	}
	addr := ":" + port
	log.Printf("listening on %s", addr) // Output: 2006/01/02 15:04:05 listening on :18080 (输出可能变化/不固定：包含时间戳)
	err := http.ListenAndServe(addr, withJSONNotFound(mux))
	log.Printf("server stopped: %v", err) // Output: 2006/01/02 15:04:05 server stopped: listen tcp :18080: bind: address already in use (输出可能变化/不固定：取决于环境与错误)
}
