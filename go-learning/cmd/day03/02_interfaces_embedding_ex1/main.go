package main

import (
	"context"
	"errors"
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

var ErrInvalidQuery = errors.New("invalid query")

type UserRepo interface {
	List(ctx context.Context, q ListUsersQuery) (items []User, total int, err error)
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService { return &UserService{repo: repo} }

func (s *UserService) ListUsers(ctx context.Context, q ListUsersQuery) ([]User, int, error) {
	if q.Page < 1 || q.Size < 1 || q.Size > 100 {
		return nil, 0, fmt.Errorf("%w", ErrInvalidQuery)
	}
	return s.repo.List(ctx, q)
}

// Exercise 1 Answer:
// - 用一个“静态/假 repo（fake repo）”实现接口，让 service 在不连数据库的情况下也能工作（便于演示/测试）。
type StaticUserRepo struct{}

func (r StaticUserRepo) List(_ context.Context, q ListUsersQuery) ([]User, int, error) {
	items := []User{
		{ID: 1, Email: "demo@example.com", Name: "Demo", Role: "viewer"},
	}
	total := len(items)
	// 简化：直接按 q.Size 截断（避免引入额外概念）
	if q.Size < total {
		items = items[:q.Size]
	}
	return items, total, nil
}

func main() {
	fmt.Println("== Day03.2 ex1: fake repo for service ==") // Output: == Day03.2 ex1: fake repo for service ==

	svc := NewUserService(StaticUserRepo{})
	items, total, err := svc.ListUsers(context.Background(), ListUsersQuery{Page: 1, Size: 1})
	fmt.Printf("total=%d err=%v\n", total, err) // Output: total=1 err=<nil>
	for _, u := range items {
		fmt.Printf("user: id=%d email=%s name=%s role=%s\n", u.ID, u.Email, u.Name, u.Role) // Output: user: id=1 email=demo@example.com name=Demo role=viewer
	}
}

