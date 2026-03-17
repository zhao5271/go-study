package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"example.com/go-learning/internal/httpkit"
)

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

var (
	errDB           = errors.New("db error")
)

func onlyGET(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			httpkit.WriteError(w, http.StatusMethodNotAllowed, 10001, "METHOD_NOT_ALLOWED") // Output: {"code":10001,"message":"METHOD_NOT_ALLOWED"}\n
			return
		}
		next(w, r)
	}
}

func withJSONNotFound(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, pattern := mux.Handler(r)
		if pattern == "" || h == nil {
			httpkit.WriteError(w, http.StatusNotFound, 10004, "NOT_FOUND") // Output: {"code":10004,"message":"NOT_FOUND"}\n
			return
		}
		h.ServeHTTP(w, r)
	})
}

func listUsers(ctx context.Context, db *sql.DB, page int, size int, search string) (items []User, total int, err error) {
	whereSQL := ""
	args := []interface{}{}
	search = strings.TrimSpace(search)
	if search != "" {
		whereSQL = "WHERE email LIKE ? OR name LIKE ?"
		like := "%" + search + "%"
		args = append(args, like, like)
	}

	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: count: %v", errDB, err)
	}

	offset := (page - 1) * size
	query := "SELECT id, email, name, role FROM users " + whereSQL + " ORDER BY id LIMIT ? OFFSET ?"
	args = append(args, size, offset)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: list: %v", errDB, err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Role); err != nil {
			return nil, 0, fmt.Errorf("%w: scan: %v", errDB, err)
		}
		items = append(items, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%w: rows: %v", errDB, err)
	}
	return items, total, nil
}

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "app:app@tcp(127.0.0.1:3307)/go_admin?parseTime=true"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("open db err=%v", err) // Output: open db err=<nil> (输出可能变化/不固定：取决于 dsn)
		return
	}
	defer func() { _ = db.Close() }()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", onlyGET(func(w http.ResponseWriter, r *http.Request) {
		httpkit.WriteJSON(w, http.StatusOK, httpkit.APIResponse{Code: 0, Message: "OK", Data: map[string]bool{"ok": true}}) // Output: {"code":0,"message":"OK","data":{"ok":true}}\n
	}))

	mux.HandleFunc("/api/v1/users", onlyGET(func(w http.ResponseWriter, r *http.Request) {
		page, size, err := httpkit.ParsePageSize(r)
		if err != nil {
			httpkit.WriteError(w, http.StatusBadRequest, 10002, "INVALID_QUERY") // Output: {"code":10002,"message":"INVALID_QUERY"}\n
			return
		}
		search := r.URL.Query().Get("search")

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		items, total, err := listUsers(ctx, db, page, size, search)
		if err != nil {
			if errors.Is(err, errDB) {
				httpkit.WriteError(w, http.StatusInternalServerError, 20001, "DB_ERROR") // Output: {"code":20001,"message":"DB_ERROR"}\n
				return
			}
			httpkit.WriteError(w, http.StatusInternalServerError, 20002, "INTERNAL_ERROR") // Output: {"code":20002,"message":"INTERNAL_ERROR"}\n
			return
		}

		log.Printf("list_users page=%d size=%d search=%q total=%d", page, size, search, total) // Output: 2006/01/02 15:04:05 list_users page=1 size=2 search="a" total=2 (输出可能变化/不固定：包含时间戳/数据变化)
		httpkit.WriteJSON(w, http.StatusOK, httpkit.APIResponse{Code: 0, Message: "OK", Data: ListUsersData{Items: items, Page: page, Size: size, Total: total}}) // Output: {"code":0,"message":"OK","data":{"items":[...],"page":1,"size":2,"total":2}}\n
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "18080"
	}
	addr := ":" + port
	log.Printf("listening on %s", addr) // Output: 2006/01/02 15:04:05 listening on :18080 (输出可能变化/不固定：包含时间戳)
	err = http.ListenAndServe(addr, withJSONNotFound(mux))
	log.Printf("server stopped: %v", err) // Output: 2006/01/02 15:04:05 server stopped: listen tcp :18080: bind: address already in use (输出可能变化/不固定：取决于环境与错误)
}
