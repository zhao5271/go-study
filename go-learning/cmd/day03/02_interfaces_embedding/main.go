package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type User struct {
	ID    int64
	Email string
	Name  string
	Role  string
}

type ListUsersQuery struct {
	Page   int
	Size   int
	Search string
}

var ErrInvalidQuery = errors.New("invalid query")

// UserRepo 是“数据访问边界”：service 不关心数据来自内存/数据库/HTTP，只关心它能否按语义取到数据。
type UserRepo interface {
	List(ctx context.Context, q ListUsersQuery) (items []User, total int, err error)
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ListUsers(ctx context.Context, q ListUsersQuery) ([]User, int, error) {
	if q.Page < 1 || q.Size < 1 || q.Size > 100 {
		return nil, 0, fmt.Errorf("%w: page must be >=1, size must be in [1,100]", ErrInvalidQuery)
	}
	return s.repo.List(ctx, q)
}

// MemoryUserRepo：一个“可运行、可演示”的 repo 实现（后面换成 MySQL repo 不影响 service）。
type MemoryUserRepo struct {
	users []User
}

func NewMemoryUserRepo(users []User) *MemoryUserRepo {
	return &MemoryUserRepo{users: users}
}

func (r *MemoryUserRepo) List(_ context.Context, q ListUsersQuery) ([]User, int, error) {
	search := strings.TrimSpace(strings.ToLower(q.Search))
	filtered := make([]User, 0, len(r.users))
	for _, u := range r.users {
		if search == "" {
			filtered = append(filtered, u)
			continue
		}
		if strings.Contains(strings.ToLower(u.Email), search) || strings.Contains(strings.ToLower(u.Name), search) {
			filtered = append(filtered, u)
		}
	}

	total := len(filtered)
	start := (q.Page - 1) * q.Size
	if start > total {
		start = total
	}
	end := start + q.Size
	if end > total {
		end = total
	}
	return filtered[start:end], total, nil
}

// LoggingUserRepo：用 embedding 做“装饰器（decorator）”，在不改业务逻辑的情况下加日志/统计/权限等横切逻辑。
type LoggingUserRepo struct {
	UserRepo // embedding：把被包装的 repo “塞进来”，作为默认实现
}

func (r LoggingUserRepo) List(ctx context.Context, q ListUsersQuery) ([]User, int, error) {
	fmt.Printf("[repo] List called page=%d size=%d search=%q\n", q.Page, q.Size, q.Search) // Output: [repo] List called page=1 size=2 search="a"
	items, total, err := r.UserRepo.List(ctx, q)
	fmt.Printf("[repo] List done total=%d err=%v\n", total, err) // Output: [repo] List done total=2 err=<nil>
	return items, total, err
}

func main() {
	fmt.Println("== Day03.2: interface + implicit impl + embedding (repo/service layering) ==") // Output: == Day03.2: interface + implicit impl + embedding (repo/service layering) ==

	mem := NewMemoryUserRepo([]User{
		{ID: 1, Email: "admin@corp.test", Name: "Admin", Role: "admin"},
		{ID: 2, Email: "alice@corp.test", Name: "Alice", Role: "editor"},
		{ID: 3, Email: "bob@corp.test", Name: "Bob", Role: "viewer"},
	})

	// service 只依赖接口；repo 的实现可替换（内存/数据库/mock），这是工程分层的关键。
	repo := LoggingUserRepo{UserRepo: mem}
	svc := NewUserService(repo)

	items, total, err := svc.ListUsers(context.Background(), ListUsersQuery{Page: 1, Size: 2, Search: "a"})
	fmt.Printf("svc.ListUsers total=%d err=%v\n", total, err) // Output: svc.ListUsers total=2 err=<nil>
	for _, u := range items {
		fmt.Printf("user: id=%d email=%s name=%s role=%s\n", u.ID, u.Email, u.Name, u.Role) // Output: user: id=1 email=admin@example.com name=Admin role=admin
	}

	_, _, err = svc.ListUsers(context.Background(), ListUsersQuery{Page: 0, Size: 2})
	fmt.Printf("invalid query errors.Is=%v err=%v\n", errors.Is(err, ErrInvalidQuery), err) // Output: invalid query errors.Is=true err=invalid query: page must be >=1, size must be in [1,100]
}
