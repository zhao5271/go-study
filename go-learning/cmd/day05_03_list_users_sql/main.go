package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        int64
	Email     string
	Name      string
	Role      string
	CreatedAt time.Time
}

func ListUsers(ctx context.Context, db *sql.DB, page int, size int, search string) (items []User, total int, err error) {
	if page < 1 {
		return nil, 0, fmt.Errorf("page must be >= 1")
	}
	if size < 1 || size > 100 {
		return nil, 0, fmt.Errorf("size must be in [1, 100]")
	}

	whereSQL := ""
	args := []interface{}{}
	search = strings.TrimSpace(search)
	if search != "" {
		// 最小示例：用 LIKE 做模糊搜索（注意：前后缀通配会影响索引；后面会讲更好的方案）
		whereSQL = "WHERE email LIKE ? OR name LIKE ?"
		like := "%" + search + "%"
		args = append(args, like, like)
	}

	countSQL := "SELECT COUNT(*) FROM users " + whereSQL
	if err := db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	listSQL := "SELECT id, email, name, role, created_at FROM users " + whereSQL + " ORDER BY id LIMIT ? OFFSET ?"
	listArgs := append(args, size, offset)
	rows, err := db.QueryContext(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &u.CreatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	items, total, err := ListUsers(ctx, db, 1, 2, "a")
	if err != nil {
		log.Printf("ListUsers err=%v", err) // Output: ListUsers err=<nil> (输出可能变化/不固定：取决于数据库状态/入参)
		return
	}

	fmt.Printf("total=%d\n", total) // Output: total=2 (输出可能变化/不固定：取决于数据与 search)
	for _, u := range items {
		fmt.Printf("user: id=%d email=%s name=%s role=%s\n", u.ID, u.Email, u.Name, u.Role) // Output: user: id=1 email=admin@example.com name=Admin role=admin
	}
}

