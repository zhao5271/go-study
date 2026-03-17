# Go 全栈学习笔记 - Day 02：函数与错误处理（TS/Node → Go）

> 今日目标：把“函数 + (value, error) + errors.Is/As + defer + table-driven tests”变成肌肉记忆。对你这种 TS/Vue/Node 背景来说，这是从 try/catch 思维切到 Go 工程思维的关键一步。

## 1) 知识讲解：概念 → 为什么（设计动机/取舍）

### 1.1 函数签名与多返回值（为什么 Go 不推崇重载/默认参数）
**概念**
- Go 函数可以返回多个值：`func f() (a, b int)`
- 典型工程约定：`(value, error)`，把“失败”显式建模成返回值，而不是异常。

**为什么**
- 取舍：Go 牺牲了语言层面的一些花活（重载、默认参数、异常），换来可读性与工程一致性。
- 多返回值让你在不引入额外对象/结构的情况下表达更多结果（对照 TS/Node：你常用对象/元组返回）。

**对照 TS/Node**
- Node 常见：返回 `{data, error}` 或 `throw` + `try/catch`。
- Go 常见：`return data, nil` / `return zero, err`，并要求你在调用处显式处理 `err`。

#### 可运行例子（紧跟）
运行：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/day02_01_returns
```

代码（可直接阅读）：
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

---

### 1.2 错误是值：`(value, error)` + wrap + `errors.Is/As`
**概念**
- `error` 是接口：任何实现 `Error() string` 的类型都能当 error。
- **sentinel error**：用一个全局 `var ErrXxx = errors.New("...")` 表示可匹配的错误条件。
- **typed error**：用自定义类型承载结构化信息（可用 `errors.As` 提取）。
- **wrap**：用 `fmt.Errorf("...: %w", err)` 在不丢失根因的情况下加上下文；上层用 `errors.Is/As` 匹配错误链。

**为什么**
- `errors.Is/As` 让你“既能给人看的上下文”，又能“给程序看的可匹配语义”。
- 是否 `%w` 很关键：wrap 了就等于把底层错误语义暴露成 API 的一部分（对封装/抽象有影响）。

**对照 TS/Node**
- TS 里通常 `throw new Error(...)`，上层靠 `instanceof` 或 message 判断（脆弱）。
- Go 倾向于：稳定语义（sentinel/type）+ 上下文（wrap），且调用处显式处理。

#### 可运行例子（紧跟）
运行：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/day02_02_errors
```

代码（可直接阅读）：
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
	fmt.Printf("hidden=%v\n", hidden)                         // Output: hidden=hide: user not found
	fmt.Printf("errors.Unwrap(hidden)=%v\n", errors.Unwrap(hidden)) // Output: errors.Unwrap(hidden)=<nil>

	_, err = users.FindUserTyped(2)
	fmt.Printf("typed err=%v\n", err) // Output: typed err=find user 2: user 2 not found

	var nf *users.NotFoundError
	fmt.Printf("errors.As(err, *NotFoundError)=%v\n", errors.As(err, &nf)) // Output: errors.As(err, *NotFoundError)=true
	fmt.Printf("nf.Resource=%s nf.ID=%d\n", nf.Resource, nf.ID)            // Output: nf.Resource=user nf.ID=2
}
```

---

### 1.3 `if err := ...; err != nil` 与 `:=` shadowing
**概念**
- `if v, err := f(); err != nil { ... }` 是 Go 常见写法：把临时变量作用域收紧在 if 内。
- 坑：你以为更新了外层 `err`，实际上 `:=` 在 if 内创建了新变量，外层 `err` 不变。

**为什么**
- 取舍：更短、更局部、更不容易把临时变量泄漏到函数其它部分；但你必须理解 shadowing。

#### 可运行例子（紧跟）
运行：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/day02_03_shadowing
```

代码（可直接阅读）：
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

---

### 1.4 `defer`：资源释放、LIFO、参数求值时机、具名返回值
**概念**
- `defer` 常用于“打开资源后立刻写关闭”，避免遗漏。
- 规则（一定要背下来）：
  1) defer 的参数**在 defer 语句执行时就求值**
  2) defer **LIFO**（后进先出）
  3) defer 可以读写**具名返回值**

#### 可运行例子（紧跟）
运行：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/day02_04_defer
```

代码（可直接阅读）：
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

---

### 1.5 panic vs error + recover（边界兜底）
**概念**
- **业务失败**：返回 `error`（别把 panic 当 throw 用）
- **panic**：程序员错误/不可恢复情况（越界、空指针等）
- **recover**：只能在 defer 里生效，用于边界兜底（注意：只对当前 goroutine）

#### 可运行例子（紧跟）
运行：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/day02_05_panic_recover
```

代码（可直接阅读）：
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
	fmt.Println("panic demo start") // Output: panic demo start
	fmt.Printf("recovered=%v\n", safe(func() { panic("boom") })) // Output: recovered=boom
	fmt.Printf("recovered(no panic)=%v\n", safe(func() {}))      // Output: recovered(no panic)=<nil>
	fmt.Println("still running")                                // Output: still running
}

func safe(fn func()) (recovered any) {
	defer func() {
		recovered = recover()
	}()
	fn()
	return nil
}
```

---

### 1.6 Table-driven tests：把测试当第一等公民
**概念**
- Go 社区非常推崇 table-driven tests：把“测试数据表”和“断言逻辑”分离，便于扩展用例。
- 常用模式：`tests := []struct{...}{...}; for _, tt := range tests { t.Run(tt.name, func(t *testing.T){...}) }`

#### 可运行例子（紧跟）
本节暂时**不要求你运行 `go test`**（你已说明目前不需要）。先把写法与结构看懂即可：
- 测试代码在：`/Users/zhang/Desktop/go-study/codex/go-learning/internal/day02/users/users_test.go:1`
- 你想跑的时候再用：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go test ./... -run TestFindUser
```

## 2) 常见坑（结合 TS/Node 习惯对照）
- **不要用 `==` 匹配错误**：上层要用 `errors.Is/As`（尤其当你 wrap 了错误时）。
- **滥用 `%w` 会泄漏实现细节**：wrap 会把底层错误语义变成 API 承诺；不想承诺就用 `%v` 仅给人看。
- **shadowing**：`if _, err := ...` 不会更新外层 `err`。
- **defer 在循环里**：循环里 defer 可能导致资源直到函数结束才释放（文件句柄/锁）。
- **recover 只能在 defer 生效**：而且只对当前 goroutine。

## 3) 工程用法/最佳实践（真实 API 项目中怎么落地）
- API handler 的套路：`v, err := svc.Do(ctx, req); if err != nil { ... }`（错误向上返回，入口统一转换成 HTTP 响应）。
- 错误对外暴露：对外尽量用稳定语义（sentinel/type），不要把底层库错误当 API 细节泄漏。
- 错误信息写法：上下文要短（避免“failed to ... failed to ...”层层堆叠）。
- 测试优先覆盖错误分支：尤其是“not found / invalid input / timeout”等业务关键路径。

## 4) 练习策略（练习可直接作为“运用示例”，提供完整参考实现）
完整参考实现就是今天这些可运行示例 + 可运行测试：
- `go run ./cmd/day02_01_returns`
- `go run ./cmd/day02_02_errors`
- `go run ./cmd/day02_03_shadowing`
- `go run ./cmd/day02_04_defer`
- `go run ./cmd/day02_05_panic_recover`
- （可选）`go test ./... -run TestFindUser`

建议你做 2 个“改造练习”（直接改现有代码即可）：
1) 在 `FindUserSentinel` 的 not found 分支再 wrap 一层上下文（多一层 `%w`），观察 `errors.Is` 仍然为 true。
2) 在 `day02_03_shadowing` 里故意写一个“需要外层 err 更新”的逻辑分支，然后用 `=` 修复（体会 shadowing 的风险）。

## References
- 官方：`https://go.dev/blog/go1.13-errors` — `%w`、wrap、`errors.Is/As/Unwrap` 的官方语义与取舍
- 官方：`https://pkg.go.dev/errors` — `errors.Is/As/Unwrap` API 文档
- 官方：`https://pkg.go.dev/fmt#Errorf` — `fmt.Errorf` 中 `%w` 的行为说明
- 官方：`https://go.dev/blog/defer-panic-and-recover` — defer/panic/recover 的三条规则与工作机制
- 官方：`https://go.dev/blog/subtests` — `t.Run` 子测试与 table-driven 的典型用法
- 官方：`https://pkg.go.dev/testing` — testing 包与 `go test` 行为说明
- 社区：`https://github.com/uber-go/guide/blob/master/style.md` — 工程实践：错误命名、wrap 取舍、错误处理模式
