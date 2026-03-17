package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID    int64
	Email string
	Name  string
	Role  string
}

func main() {
	// 典型 DSN（本仓库的 compose 端口映射是 3307:3306）
	// app:app@tcp(127.0.0.1:3307)/go_admin?parseTime=true
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "app:app@tcp(127.0.0.1:3307)/go_admin?parseTime=true"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("open db err=%v", err) // Output: open db err=<nil> (输出可能变化/不固定：取决于 dsn)
		return
	}
	defer func() {
		_ = db.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Printf("ping err=%v", err) // Output: ping err=dial tcp 127.0.0.1:3307: connect: connection refused (输出可能变化/不固定：取决于环境/容器状态)
		return
	}
	fmt.Println("db ping ok") // Output: db ping ok

	email := fmt.Sprintf("demo_%d@example.com", time.Now().UnixNano())
	name := "Demo User"
	role := "viewer"

	res, err := db.ExecContext(ctx, `INSERT INTO users (email, name, role) VALUES (?, ?, ?)`, email, name, role)
	if err != nil {
		log.Printf("insert err=%v", err) // Output: insert err=<nil> (输出可能变化/不固定：取决于数据库约束/状态)
		return
	}
	id, _ := res.LastInsertId()
	fmt.Printf("inserted user id=%d\n", id) // Output: inserted user id=4 (输出可能变化/不固定：自增 id 取决于现有数据)

	rows, err := db.QueryContext(ctx, `SELECT id, email, name, role FROM users ORDER BY id LIMIT 5`)
	if err != nil {
		log.Printf("query err=%v", err) // Output: query err=<nil> (输出可能变化/不固定：取决于数据库状态)
		return
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var u User
		if scanErr := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Role); scanErr != nil {
			log.Printf("scan err=%v", scanErr) // Output: scan err=<nil> (输出可能变化/不固定：取决于数据类型/列)
			return
		}
		fmt.Printf("user: id=%d email=%s name=%s role=%s\n", u.ID, u.Email, u.Name, u.Role) // Output: user: id=1 email=admin@example.com name=Admin role=admin
	}
	if err := rows.Err(); err != nil {
		log.Printf("rows err=%v", err) // Output: rows err=<nil> (输出可能变化/不固定：取决于网络/驱动)
	}
}

