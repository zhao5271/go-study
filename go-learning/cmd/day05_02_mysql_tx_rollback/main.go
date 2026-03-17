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

	if err := db.PingContext(ctx); err != nil {
		log.Printf("ping err=%v", err) // Output: ping err=dial tcp 127.0.0.1:3307: connect: connection refused (输出可能变化/不固定：取决于环境/容器状态)
		return
	}
	fmt.Println("db ping ok") // Output: db ping ok

	// 用一个相对固定的 email，便于你重复运行验证事务回滚不会“半成功”
	email := "tx_demo_should_rollback@example.com"

	// 先清理遗留（非必要，但让输出更稳定）
	_, _ = db.ExecContext(ctx, `DELETE FROM audit_logs WHERE target_type='user' AND target_id IN (SELECT id FROM users WHERE email=?)`, email)
	_, _ = db.ExecContext(ctx, `DELETE FROM users WHERE email=?`, email)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("begin tx err=%v", err) // Output: begin tx err=<nil> (输出可能变化/不固定：取决于数据库状态)
		return
	}

	// 1) 插入一个用户
	res, err := tx.ExecContext(ctx, `INSERT INTO users (email, name, role) VALUES (?, ?, ?)`, email, "Tx Demo", "viewer")
	if err != nil {
		_ = tx.Rollback()
		log.Printf("insert user err=%v", err) // Output: insert user err=<nil> (输出可能变化/不固定：取决于约束/状态)
		return
	}
	userID, _ := res.LastInsertId()

	// 2) 插入一条审计日志（同一个事务里）
	_, err = tx.ExecContext(ctx, `INSERT INTO audit_logs (actor_user_id, action, target_type, target_id) VALUES (?, ?, ?, ?)`,
		1, "create_user", "user", userID)
	if err != nil {
		_ = tx.Rollback()
		log.Printf("insert audit err=%v", err) // Output: insert audit err=<nil> (输出可能变化/不固定：取决于约束/状态)
		return
	}

	// 3) 故意制造一个错误：再次插入同 email（触发 UNIQUE 约束失败）
	_, err = tx.ExecContext(ctx, `INSERT INTO users (email, name, role) VALUES (?, ?, ?)`, email, "Duplicate", "viewer")
	if err == nil {
		_ = tx.Rollback()
		log.Printf("expected duplicate error, got nil") // Output: expected duplicate error, got nil
		return
	}
	log.Printf("expected err=%v", err) // Output: expected err=Error 1062 (23000): Duplicate entry ... (输出可能变化/不固定：错误信息随驱动/版本变化)

	// 4) 回滚：应当不会留下“用户已写入，但审计日志没写入/或反之”的半成功数据
	if rbErr := tx.Rollback(); rbErr != nil {
		log.Printf("rollback err=%v", rbErr) // Output: rollback err=<nil> (输出可能变化/不固定：取决于状态)
		return
	}
	fmt.Println("rollback ok") // Output: rollback ok

	var userCount int
	_ = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE email=?`, email).Scan(&userCount)
	fmt.Printf("users count=%d\n", userCount) // Output: users count=0

	var auditCount int
	_ = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_logs WHERE target_type='user' AND target_id=?`, userID).Scan(&auditCount)
	fmt.Printf("audit_logs count=%d\n", auditCount) // Output: audit_logs count=0
}

