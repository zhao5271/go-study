# Day 03.2：接口（interface）+ 隐式实现 + embedding（工程分层：Repo/Service）

> 贯穿项目：后台管理 API（RBAC + 登录鉴权 + 列表分页检索 + CRUD + 审计日志 + Docker 部署）  
> 本节目标：用 **interface 把“业务逻辑”和“数据访问”解耦**，并用 embedding 做装饰器（日志/统计）——为后续 MySQL repo、Gin handler 分层做铺垫。

统一运行目录：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
```

---

## 知识点 1：`interface` + 隐式实现（Repo/Service 的分层边界）

### B. 一句话定义
接口（interface）是“一组方法的契约”；Go 用“方法集合匹配”来 **隐式实现** 接口，不需要 `implements` 关键字。

### C. 为什么重要（不做会怎样）
后台管理项目里 service 层如果直接依赖 MySQL 细节，你很难：
- 替换实现（内存→MySQL→Mock）
- 做可演示/可测试的最小闭环（作品集交付很吃这个）
- 在 handler/service/repo 分层时保持边界清晰

### D. 重难点拆解（2–4 条）
1) **接口要“贴业务语义”**：比如 `ListUsers(...)`，不要把 SQL 细节塞进接口里。  
2) **接口放在“使用者一侧”更合理**：通常 service 定义 `UserRepo`，repo 去实现它。  
3) **返回值形态要工程化**：列表常用 `items + total + error`，并做 page/size 边界校验（否则接口不可控）。

### E. 业务场景落地（后台管理 API）
用户列表页：service 只关心“按查询条件返回用户列表”，repo 可以是：
- 内存实现（演示/开发期）
- MySQL 实现（生产）
- fake/mock（测试）

### F. 代码示例（最小可运行）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day03_02_interfaces_embedding/main.go`
```go
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
```

### G. 怎么运行（命令 + 预期现象）
```bash
go run ./cmd/day03_02_interfaces_embedding
# Output: == Day03.2: interface + implicit impl + embedding (repo/service layering) ==
# Output: [repo] List called page=1 size=2 search="a"
# Output: [repo] List done total=2 err=<nil>
# Output: svc.ListUsers total=2 err=<nil>
# Output: user: id=1 email=admin@corp.test name=Admin role=admin
# Output: user: id=2 email=alice@corp.test name=Alice role=editor
# Output: invalid query errors.Is=true err=invalid query: page must be >=1, size must be in [1,100]
```

### H. 练习题（1–3 题，覆盖边界条件）
练习 1：实现一个 `StaticUserRepo`（不连 DB）并注入到 `UserService` 里跑通列表  
- 验收标准：
  - `go run` 能打印 total 和至少 1 个 user
  - service 层代码不需要改动（只替换 repo 实现）

### I. 参考答案（紧跟练习题）
参考答案 1（可运行）：
- 文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day03_02_interfaces_embedding_ex1/main.go`
- 运行：
```bash
go run ./cmd/day03_02_interfaces_embedding_ex1
# Output: == Day03.2 ex1: fake repo for service ==
# Output: total=1 err=<nil>
# Output: user: id=1 email=demo@example.com name=Demo role=viewer
```

---

## 知识点 2：embedding（组合）做“装饰器”（日志/统计/鉴权）

### B. 一句话定义
embedding 是把一个类型匿名嵌入到 struct 里，让它的字段/方法像“自己的一样”被提升（promoted），常用于组合与装饰器。

### C. 为什么重要（不做会怎样）
后台管理 API 的横切能力（日志、统计、权限、缓存）如果散落在各处，会变成“改一处漏一处”；用装饰器可以在不改业务逻辑的情况下加能力。

### D. 重难点拆解（2–4 条）
1) **embedding 不等于继承**：它是组合（composition），强调“有一个（has-a）”。  
2) **装饰器 = 包一层接口**：常见形态 `type X struct { SomeInterface }`。  
3) **避免过度魔法**：装饰器只做横切（log/metrics/trace），不要把业务逻辑塞进去。

### E. 业务场景落地（后台管理 API）
你后面会给 repo/service/handler 加：
- request log（请求/耗时）
- metrics（接口调用次数/错误率）
- trace（链路追踪）
用装饰器可以统一加，不用每个函数手写。

### F. 代码示例（最小可运行）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day03_02_interfaces_embedding_ex2/main.go`

### G. 怎么运行（命令 + 预期现象）
```bash
go run ./cmd/day03_02_interfaces_embedding_ex2
# Output: == Day03.2 ex2: embedding decorator for metrics ==
# Output: List called times=2
```

### H. 练习题（1–3 题，覆盖边界条件）
练习 1：把 `MetricsUserRepo` 改成同时统计 `List` 的成功/失败次数（你可以人为制造一次失败）  
- 验收标准：
  - 输出里能看到成功次数与失败次数
  - 不修改 `MemoryUserRepo`（只改装饰器）

### I. 参考答案
参考答案 1（思路）：  
在装饰器里 `items,total,err := r.UserRepo.List(...)` 后：
- `err == nil` → success++
- `err != nil` → fail++
然后打印统计（打印要写典型输出注释）。

---

## References
- 官方：Effective Go（Interfaces）https://go.dev/doc/effective_go
- 官方：Go Spec（Method sets / Interface types）https://go.dev/ref/spec
