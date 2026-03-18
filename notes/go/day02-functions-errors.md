# Day 02：函数与错误处理（把失败建模成可控的工程接口）

> 贯穿项目：后台管理 API（RBAC + 登录鉴权 + 列表分页检索 + CRUD + 审计日志 + Docker 部署）  
> 本日目标：把 `(value, error)`、wrap、`errors.Is/As`、defer 变成肌肉记忆（后面做 HTTP/MySQL 都靠它）。

统一运行目录：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
```

---

## 知识点 1：多返回值 + 命名返回值（Go 的“函数签名表达力”）

### B. 一句话定义
Go 函数可以返回多个值；命名返回值会变成真实变量（默认是零值）。

### C. 为什么重要（不做会怎样）
你后面写 service/repo 时，常见签名就是 `(...)(T, error)`；如果你对多返回值不熟，代码会变得啰嗦且容易漏错误处理。

### D. 重难点拆解（2–4 条）
1) **命名返回值不是“语法糖”**：它是真变量，会有零值。  
2) **少用“裸返回”**：只在非常短小函数里使用；否则可读性下降（你自己未来会骂自己）。  
3) **工程约定**：对外暴露接口优先 `(T, error)`，不要搞魔法返回。

### E. 业务场景落地
例如“创建用户”：返回 `userID, err`；“列表分页”：返回 `items, total, err`。

### F. 代码示例（最小可运行）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day02/01_returns/main.go`
```go
package main

import "fmt"

func main() {
	fmt.Println("== Day02.1: functions + multiple returns + named returns ==") // Output: == Day02.1: functions + multiple returns + named returns ==

	x, y := split(10)
	fmt.Printf("split(10) => x=%d y=%d\n", x, y) // Output: split(10) => x=4 y=6

	q, r := divmod(17, 5)
	fmt.Printf("divmod(17,5) => q=%d r=%d\n", q, r) // Output: divmod(17,5) => q=3 r=2

	fmt.Printf("namedZero()=%d\n", namedZero()) // Output: namedZero()=0
}

// split demonstrates named return values. They are real variables (zero-valued).
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}

func divmod(a, b int) (int, int) {
	return a / b, a % b
}

func namedZero() (x int) {
	return
}
```

### G. 怎么运行（命令 + 预期现象）
```bash
go run ./cmd/day02/01_returns
# Output: == Day02.1: functions + multiple returns + named returns ==
# Output: split(10) => x=4 y=6
# Output: divmod(17,5) => q=3 r=2
# Output: namedZero()=0
```

### H. 练习题（1–3 题）
练习 1：写一个 `ListUsers(page,size)` 的函数签名（先不实现），返回 `items, total, err`  
- 验收标准：签名是 `([]T, int, error)` 或 `([]User, int, error)` 这种形态，不要把 error 藏到结构体里

### I. 参考答案
参考答案 1（示例签名）：
```go
func ListUsers(page, size int) (items []User, total int, err error) { /* ... */ }
```

---

## 知识点 2：错误是值：sentinel / typed error + wrap（`%w`）+ `errors.Is/As`

### B. 一句话定义
Go 的错误是普通值：你可以用“可匹配语义（sentinel/type）+ 上下文（wrap）”同时满足“给程序判断”和“给人排障”。

### C. 为什么重要（不做会怎样）
后台管理 API 的错误要稳定：前端要按错误码处理，后端要按错误语义做分支；如果你只靠字符串 message，系统会非常脆弱。

### D. 重难点拆解（2–4 条）
1) **不要用 `==` 判断错误**：一旦 wrap 了就不可靠，应该 `errors.Is/As`。  
2) **`%w` 是“对上层承诺语义”**：你 wrap 了，就等于允许上层用 `errors.Is/As` 看见底层语义；不想暴露就用 `%v`。  
3) **sentinel vs typed**：sentinel 更简单；typed 能携带结构化字段（Resource/ID），更适合做统一错误映射。

### E. 业务场景落地
例如 repo 层“找不到用户”：
- service 层要能判断是 NotFound（映射 404）还是 DB 真的挂了（映射 500）

### F. 代码示例（最小可运行）
入口文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day02/02_errors/main.go`
```go
package main

import (
	"errors"
	"fmt"

	"example.com/go-learning/internal/day02/users"
)

func main() {
	fmt.Println("== Day02.2: (value, error) + wrap + errors.Is/As ==") // Output: == Day02.2: (value, error) + wrap + errors.Is/As ==

	_, err := users.FindUserSentinel(2)
	fmt.Printf("sentinel err=%v\n", err) // Output: sentinel err=find user 2: user not found

	fmt.Printf("errors.Is(err, ErrUserNotFound)=%v\n", errors.Is(err, users.ErrUserNotFound)) // Output: errors.Is(err, ErrUserNotFound)=true
	fmt.Printf("errors.Unwrap(err)=%v\n", errors.Unwrap(err))                                 // Output: errors.Unwrap(err)=user not found

	hidden := fmt.Errorf("hide: %v", users.ErrUserNotFound)
	fmt.Printf("hidden=%v\n", hidden)                               // Output: hidden=hide: user not found
	fmt.Printf("errors.Unwrap(hidden)=%v\n", errors.Unwrap(hidden)) // Output: errors.Unwrap(hidden)=<nil>

	_, err = users.FindUserTyped(2)
	fmt.Printf("typed err=%v\n", err) // Output: typed err=find user 2: user 2 not found

	var nf *users.NotFoundError
	fmt.Printf("errors.As(err, *NotFoundError)=%v\n", errors.As(err, &nf)) // Output: errors.As(err, *NotFoundError)=true
	fmt.Printf("nf.Resource=%s nf.ID=%d\n", nf.Resource, nf.ID)            // Output: nf.Resource=user nf.ID=2
}
```

被调用的内部实现（为了你阅读完整链路）：

文件：`/Users/zhang/Desktop/go-study/codex/go-learning/internal/day02/users/users.go`
```go
package users

import (
	"errors"
	"fmt"
)

type User struct {
	ID   int
	Name string
}

// ErrUserNotFound is a sentinel error. Callers should match it with errors.Is
// (not ==), because we return it wrapped with context.
var ErrUserNotFound = errors.New("user not found")

// NotFoundError is a typed error. Callers can inspect it with errors.As.
type NotFoundError struct {
	Resource string
	ID       int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s %d not found", e.Resource, e.ID)
}

func FindUserSentinel(id int) (User, error) {
	if id == 1 {
		return User{ID: 1, Name: "Gopher"}, nil
	}
	return User{}, fmt.Errorf("find user %d: %w", id, ErrUserNotFound)
}

func FindUserTyped(id int) (User, error) {
	if id == 1 {
		return User{ID: 1, Name: "Gopher"}, nil
	}
	return User{}, fmt.Errorf("find user %d: %w", id, &NotFoundError{Resource: "user", ID: id})
}
```

### G. 怎么运行（命令 + 预期现象）
```bash
go run ./cmd/day02/02_errors
# Output: == Day02.2: (value, error) + wrap + errors.Is/As ==
# Output: sentinel err=find user 2: user not found
# Output: errors.Is(err, ErrUserNotFound)=true
# Output: errors.Unwrap(err)=user not found
# Output: hidden=hide: user not found
# Output: errors.Unwrap(hidden)=<nil>
# Output: typed err=find user 2: user 2 not found
# Output: errors.As(err, *NotFoundError)=true
# Output: nf.Resource=user nf.ID=2
```

### H. 练习题（1–3 题）
练习 1：把 `FindUserSentinel(2)` 的错误再 wrap 一层上下文（多包一层 `%w`），验证 `errors.Is` 仍然是 true  
- 验收标准：`errors.Is(err, ErrUserNotFound)=true` 仍成立

练习 2：把 typed error 增加一个字段 `Reason`（字符串），并在上层打印出来  
- 验收标准：你能通过 `errors.As` 拿到 `Reason`

### I. 参考答案
参考答案 1（可运行做法）：  
在 `users.FindUserSentinel` 的 not found 分支再包一层：
```go
err := fmt.Errorf("find user %d: %w", id, ErrUserNotFound)
return User{}, fmt.Errorf("repo: %w", err)
```
然后运行 `go run ./cmd/day02/02_errors`，`errors.Is` 仍然为 true。

参考答案 2（可运行做法）：  
给 `NotFoundError` 加字段 `Reason string`，并在构造处赋值；上层 `nf.Reason` 打印（打印要写典型输出注释）。

---

## 知识点 3：控制流里的错误处理坑（shadowing）+ `defer` 三条规则 + panic/recover 的边界

### B. 一句话定义
`if err := ...` 会缩小变量作用域但可能触发 shadowing；`defer` 用来可靠释放资源；panic/recover 只做边界兜底，不做业务流程。

### C. 为什么重要（不做会怎样）
你后面写 DB/HTTP 时会大量出现：打开资源后释放、错误向上返回、入口兜底防崩；这三件事搞错，会直接变成线上事故（资源泄漏/错误被吞/服务崩溃）。

### D. 重难点拆解（2–4 条）
1) **shadowing**：`:=` 在块内可能新建同名变量，外层变量不会被更新。  
2) **defer 参数求值时机**：defer 的参数在“声明 defer 时”就求值。  
3) **panic/recover 边界**：recover 只能在 defer 内生效，且只对当前 goroutine；业务失败应返回 error。

### E. 业务场景落地
- 你写 `db.Query()` 后一定要 `defer rows.Close()`；写 handler 要有统一错误响应；服务入口可加 recover 兜底避免整个进程崩掉（但要记录日志）。

### F. 代码示例（最小可运行）
文件 1（shadowing）：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day02/03_shadowing/main.go`
```go
package main

import (
	"fmt"

	"example.com/go-learning/internal/day02/users"
)

func main() {
	fmt.Println("== Day02.3: if init + := shadowing ==") // Output: == Day02.3: if init + := shadowing ==

	user, err := users.FindUserSentinel(1)
	fmt.Printf("outer: user.ID=%d err=%v\n", user.ID, err) // Output: outer: user.ID=1 err=<nil>

	// Common pitfall: := inside if creates a NEW err that does not update outer err.
	if _, err := users.FindUserSentinel(2); err != nil {
		fmt.Printf("inside if: err=%v\n", err) // Output: inside if: err=find user 2: user not found
	}
	fmt.Printf("after if: outer err=%v\n", err) // Output: after if: outer err=<nil>

	// If you really want to update outer err, use assignment (=), not :=.
	_, err = users.FindUserSentinel(2)
	fmt.Printf("after assignment: err=%v\n", err) // Output: after assignment: err=find user 2: user not found
}
```

文件 2（defer）：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day02/04_defer/main.go`
```go
package main

import "fmt"

func main() {
	fmt.Println("== Day02.4: defer (LIFO / args eval / named return) ==") // Output: == Day02.4: defer (LIFO / args eval / named return) ==

	demoLIFO()
	demoArgsEvaluatedAtDefer()

	fmt.Printf("namedReturn()=%d\n", namedReturn()) // Output: namedReturn()=2
}

func demoLIFO() {
	fmt.Println("demoLIFO start") // Output: demoLIFO start
	defer fmt.Println("defer 1")  // Output: defer 1
	defer fmt.Println("defer 2")  // Output: defer 2
	fmt.Println("demoLIFO end")   // Output: demoLIFO end
}

func demoArgsEvaluatedAtDefer() {
	fmt.Println("demoArgs start") // Output: demoArgs start
	x := 1
	defer fmt.Printf("defer x=%d\n", x) // Output: defer x=1
	x = 2
	fmt.Printf("now x=%d\n", x) // Output: now x=2
	fmt.Println("demoArgs end") // Output: demoArgs end
}

func namedReturn() (result int) {
	defer func() {
		result++
	}()
	return 1
}
```

文件 3（panic/recover）：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day02/05_panic_recover/main.go`
```go
package main

import (
	"fmt"

	"example.com/go-learning/internal/day02/users"
)

func main() {
	fmt.Println("== Day02.5: panic vs error + recover ==") // Output: == Day02.5: panic vs error + recover ==

	// Business logic failures should return error, not panic.
	_, err := users.FindUserSentinel(2)
	fmt.Printf("business err=%v\n", err) // Output: business err=find user 2: user not found

	// Panic is for programmer errors / truly unrecoverable situations.
	fmt.Println("panic demo start")                              // Output: panic demo start
	fmt.Printf("recovered=%v\n", safe(func() { panic("boom") })) // Output: recovered=boom
	fmt.Printf("recovered(no panic)=%v\n", safe(func() {}))      // Output: recovered(no panic)=<nil>
	fmt.Println("still running")                                 // Output: still running
}

func safe(fn func()) (recovered any) {
	defer func() {
		recovered = recover()
	}()
	fn()
	return nil
}
```

### G. 怎么运行（命令 + 预期现象）
```bash
go run ./cmd/day02/03_shadowing
# Output: == Day02.3: if init + := shadowing ==
# Output: outer: user.ID=1 err=<nil>
# Output: inside if: err=find user 2: user not found
# Output: after if: outer err=<nil>
# Output: after assignment: err=find user 2: user not found

go run ./cmd/day02/04_defer
# Output: == Day02.4: defer (LIFO / args eval / named return) ==
# Output: demoLIFO start
# Output: demoLIFO end
# Output: defer 2
# Output: defer 1
# Output: demoArgs start
# Output: now x=2
# Output: demoArgs end
# Output: defer x=1
# Output: namedReturn()=2

go run ./cmd/day02/05_panic_recover
# Output: == Day02.5: panic vs error + recover ==
# Output: business err=find user 2: user not found
# Output: panic demo start
# Output: recovered=boom
# Output: recovered(no panic)=<nil>
# Output: still running
```

### H. 练习题（1–3 题）
练习 1：在 `day02/03_shadowing` 里增加一个“如果发生错误就 return”的逻辑，确保使用 `=` 更新外层 err  
- 验收标准：当找不到用户时，函数最终能拿到外层 err（不是 `<nil>`）

练习 2：写下 defer 三条规则，并用 `day02/04_defer` 的输出逐条对照证明  
- 验收标准：你能解释为什么 `defer x=1` 而不是 2

### I. 参考答案
参考答案 1：  
把 `if _, err := ...` 改成：
```go
_, err = users.FindUserSentinel(2)
if err != nil {
	// ...
}
```

参考答案 2：  
对照输出：`defer x=1` 是因为 defer 的参数在“写 defer 的那一刻”就求值了（不是执行时求值）。

---

## References
- 官方：Go 1.13 errors（wrap/%w/Is/As）https://go.dev/blog/go1.13-errors
- 官方：errors 包 https://pkg.go.dev/errors
- 官方：fmt.Errorf（%w）https://pkg.go.dev/fmt#Errorf
- 官方：defer/panic/recover https://go.dev/blog/defer-panic-and-recover
- 社区：Uber Go Style Guide（工程错误处理取舍）https://github.com/uber-go/guide/blob/master/style.md
