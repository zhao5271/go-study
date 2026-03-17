package main

import (
	"context"
	"fmt"
)

type User struct {
	ID    int64
	Email string
	Name  string
	Role  string
}

type ListUsersQuery struct {
	Page int
	Size int
}

type UserRepo interface {
	List(ctx context.Context, q ListUsersQuery) (items []User, total int, err error)
}

// Memory repo（最小实现）
type MemoryUserRepo struct {
	users []User
}

func (r *MemoryUserRepo) List(_ context.Context, q ListUsersQuery) ([]User, int, error) {
	total := len(r.users)
	start := (q.Page - 1) * q.Size
	if start > total {
		start = total
	}
	end := start + q.Size
	if end > total {
		end = total
	}
	return r.users[start:end], total, nil
}

// Exercise 2 Answer:
// - 用 embedding 写一个 metrics 装饰器：不改底层 repo，就能统计调用次数。
type MetricsUserRepo struct {
	UserRepo // embedding：默认把所有方法转发给被包装对象
	ListN    int
}

func (r *MetricsUserRepo) List(ctx context.Context, q ListUsersQuery) ([]User, int, error) {
	r.ListN++
	return r.UserRepo.List(ctx, q)
}

func main() {
	fmt.Println("== Day03.2 ex2: embedding decorator for metrics ==") // Output: == Day03.2 ex2: embedding decorator for metrics ==

	base := &MemoryUserRepo{users: []User{
		{ID: 1, Email: "a@example.com", Name: "A", Role: "viewer"},
		{ID: 2, Email: "b@example.com", Name: "B", Role: "viewer"},
	}}

	repo := &MetricsUserRepo{UserRepo: base}
	_, _, _ = repo.List(context.Background(), ListUsersQuery{Page: 1, Size: 1})
	_, _, _ = repo.List(context.Background(), ListUsersQuery{Page: 1, Size: 2})

	fmt.Printf("List called times=%d\n", repo.ListN) // Output: List called times=2
}

