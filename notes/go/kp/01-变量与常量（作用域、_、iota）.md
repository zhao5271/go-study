# 01 变量与常量（作用域、`_`、`iota`）

> 目标：把“变量/常量的可见范围 + 忽略值 + iota 常量枚举”一次讲透，并能直接跑一个最小示例。  
> 贯穿项目落地：后台管理 API 里会经常遇到“临时变量/错误变量的作用域”“忽略不需要的返回值”“用常量表达稳定的业务枚举/错误码/角色”等。

## 关联复习（已学过内容摘要，不让你跳去翻旧笔记）

1) `:=`（短变量声明）
- 只能在函数体内使用：`x := 1`
- 它做的是“声明 + 赋值 + 类型推导”

2) shadowing（遮蔽）= `:=` 在块里可能“新建同名变量”
```go
package main

import "fmt"

func main() {
	x := 1
	if true {
		x := 2
		fmt.Printf("inner x=%d\n", x) // Output: inner x=2
	}
	fmt.Printf("outer x=%d\n", x) // Output: outer x=1
}
```

3) 常量（`const`）的直觉
- `const` 表达“不会变的值”（角色/状态/错误码类别…），避免 magic number/string

## 资料来源说明
- NotebookLM 本次查询失败（`net::ERR_CONNECTION_CLOSED`），本节结论以 Go 官方文档为准（见文末 References）。

---

## 知识点 1：作用域（scope）+ 块级作用域（block scope）+ `:=` 的遮蔽

你可以把 Go 的作用域理解成：**大括号 `{}` 就是一道“变量可见边界”**。

最常见的 3 个层级：
- **包级（package scope）**：包里任何文件都能用（例如 `var ErrXxx`）
- **函数级**：函数体内有效
- **块级（block scope）**：`if/for/switch/{...}` 的大括号内有效

为什么重要（后台管理 API 场景）：
- 你会大量写 `if err != nil { ... }` / `if v, err := f(); err != nil { ... }`。  
  如果不清楚块作用域，很容易把“外层变量没更新”的 bug 带到 handler/service/repo 里。

关键坑（只记 2 条就够用）：
1) `:=` 在块里**可能新建同名变量**，外层变量不会变（这就是 shadowing）。  
2) 推荐写法 `if v, err := f(); err != nil { ... }` 的优点是**缩小变量作用域**，降低误用概率；但你必须清楚它不会影响外层同名变量。

下面的“完整可运行示例代码”里，`demoScope()` 会把这个坑打印出来，你可以直接跑看输出。

---

## 知识点 2：匿名变量 `_`（blank identifier）= 显式“我不需要这个值”

一句话：`_` 是一个“**丢弃桶**”，把你不需要的返回值/循环变量明确丢掉。

为什么重要（后台管理 API 场景）：
- 标准库/你自己的函数经常多返回值：`(value, error)`、`(rows, err)`、`(affected, err)`。  
  你只想要其中一个，就用 `_` 表达“我刻意忽略”，并通过编译器的“未使用即报错”保持代码干净。

关键坑：
1) `_` 不是“随便写写”，它会让你**永久失去那个值**；别把重要的错误/返回值用 `_` 吞了。  
2) `range` 里如果你不用 index/value，优先用 `_`，避免写一个没意义的变量名。

下面的“完整可运行示例代码”里，`demoBlankIdentifier()` 覆盖了“忽略返回值”和 “range 忽略 index” 两种常见用法。

---

## 知识点 3：`const` + `iota`（常量计数器）= 稳定枚举/位标志

一句话：`iota` 只在 `const (...)` 块里出现，**从 0 开始、每行 +1**，常用于枚举/位标志。

为什么重要（后台管理 API 场景）：
- 角色（Admin/Viewer）、状态（Enabled/Disabled）、错误码类别、审计动作类型……这些都是**稳定、可复用、可搜索**的常量，适合用 `const` 表达，而不是散落的 magic number/string。

关键坑：
1) `iota` **按“行”递增**（不是按“写了 iota 才递增”）：哪怕某行你写了显式值，它也会继续 +1。  
2) `iota` **每个 const 块都会重置为 0**：跨文件/跨块不共享计数。

下面的“完整可运行示例代码”里，`demoIota()` 给了两种最常见模式：枚举（Role）+ 位标志/容量单位（KB/MB/GB）。

---

## 完整可运行示例代码（全文在笔记里，不让你去跳转找）

> 说明：这段代码也会同步落盘到 `go-learning/cmd/kp/vars-const-scope-iota/main.go`（便于你直接 `go run` 跑起来）。

```go
package main

import "fmt"

func main() {
	fmt.Println("== KP: vars/const/scope/_/iota ==") // Output: == KP: vars/const/scope/_/iota ==

	demoScope()
	fmt.Println("----") // Output: ----

	demoBlankIdentifier()
	fmt.Println("----") // Output: ----

	demoIota()
}

func demoScope() {
	fmt.Println("[1] scope + block scope + := shadowing") // Output: [1] scope + block scope + := shadowing

	x := 1
	if true {
		x := 2
		fmt.Printf("inner x=%d\n", x) // Output: inner x=2
	}
	fmt.Printf("outer x=%d\n", x) // Output: outer x=1

	// Tip: 如果你想更新外层 x，就在 block 里用 `x = 2`（而不是 `x := 2`）。
}

func demoBlankIdentifier() {
	fmt.Println("[2] blank identifier (_)") // Output: [2] blank identifier (_)

	q, _ := divMod(10, 3)
	fmt.Printf("10/3 quotient=%d\n", q) // Output: 10/3 quotient=3

	names := []string{"alice", "bob"}
	for _, name := range names {
		fmt.Printf("name=%s\n", name) // Output: name=alice（第一次迭代）；Output: name=bob（第二次迭代）
	}
}

func divMod(a, b int) (q int, r int) {
	return a / b, a % b
}

type Role int

const (
	RoleAdmin Role = iota + 1
	RoleEditor
	RoleViewer
)

const (
	_ = iota // skip 0
	KB = 1 << (10 * iota)
	MB
	GB
)

func demoIota() {
	fmt.Println("[3] const + iota") // Output: [3] const + iota

	fmt.Printf("RoleAdmin=%d RoleEditor=%d RoleViewer=%d\n", RoleAdmin, RoleEditor, RoleViewer) // Output: RoleAdmin=1 RoleEditor=2 RoleViewer=3
	fmt.Printf("KB=%d MB=%d GB=%d\n", KB, MB, GB)                                              // Output: KB=1024 MB=1048576 GB=1073741824
}
```

---

## 怎么运行 + 预期现象（带典型输出）

```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/kp/vars-const-scope-iota
```

```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/kp/vars-const-scope-iota
# Output: == KP: vars/const/scope/_/iota ==
# Output: [1] scope + block scope + := shadowing
# Output: inner x=2
# Output: outer x=1
# Output: ----
# Output: [2] blank identifier (_)
# Output: 10/3 quotient=3
# Output: name=alice
# Output: name=bob
# Output: ----
# Output: [3] const + iota
# Output: RoleAdmin=1 RoleEditor=2 RoleViewer=3
# Output: KB=1024 MB=1048576 GB=1073741824
```

---

## 自检/练习

练习 1：修复 shadowing，让外层变量被更新  
- 要求：把 `demoScope()` 里 `if true { x := 2 ... }` 改成 `x = 2`
- 验收：输出里 `outer x=2`

练习 2：不要丢弃余数，把它打印出来  
- 要求：在 `demoBlankIdentifier()` 里把 `_` 改成变量名 `rem` 并打印它
- 验收：输出里出现 `10%3 remainder=1`

练习 3：扩展枚举常量，增加 `RoleGuest`  
- 要求：给 `Role` 增加 `RoleGuest`，并把它也打印出来
- 验收：`RoleGuest=4`（连续递增）

---

## 参考答案（每题给完整可运行代码，不写“自行修改/照抄”）

> 说明：下面每份代码都可以直接保存成 `main.go` 运行。为了避免你来回对比，我把“完整文件”都贴出来了。

### 参考答案 1：练习 1（修复 shadowing）
```go
package main

import "fmt"

func main() {
	fmt.Println("== EX1: fix shadowing ==") // Output: == EX1: fix shadowing ==
	demoScope()
}

func demoScope() {
	fmt.Println("[1] scope + block scope + := shadowing (fixed)") // Output: [1] scope + block scope + := shadowing (fixed)

	x := 1
	if true {
		x = 2
		fmt.Printf("inner x=%d\n", x) // Output: inner x=2
	}
	fmt.Printf("outer x=%d\n", x) // Output: outer x=2
}
```

### 参考答案 2：练习 2（打印余数，不再用 `_` 丢弃）
```go
package main

import "fmt"

func main() {
	fmt.Println("== EX2: print remainder ==") // Output: == EX2: print remainder ==
	demoBlankIdentifier()
}

func demoBlankIdentifier() {
	fmt.Println("[2] blank identifier (_) -> keep remainder") // Output: [2] blank identifier (_) -> keep remainder

	q, rem := divMod(10, 3)
	fmt.Printf("10/3 quotient=%d\n", q)    // Output: 10/3 quotient=3
	fmt.Printf("10%%3 remainder=%d\n", rem) // Output: 10%3 remainder=1
}

func divMod(a, b int) (q int, r int) {
	return a / b, a % b
}
```

### 参考答案 3：练习 3（增加 RoleGuest 并验证递增）
```go
package main

import "fmt"

type Role int

const (
	RoleAdmin Role = iota + 1
	RoleEditor
	RoleViewer
	RoleGuest
)

func main() {
	fmt.Println("== EX3: iota enum extend ==") // Output: == EX3: iota enum extend ==
	fmt.Printf("RoleAdmin=%d RoleEditor=%d RoleViewer=%d RoleGuest=%d\n", RoleAdmin, RoleEditor, RoleViewer, RoleGuest) // Output: RoleAdmin=1 RoleEditor=2 RoleViewer=3 RoleGuest=4
}
```

---

## References
- 官方：Go Spec（Constants / Iota / Scopes / Blank identifier）https://go.dev/ref/spec  
- 官方：Effective Go（Constants / iota 常见用法）https://go.dev/doc/effective_go  
- 官方：A Tour of Go（Constants）https://go.dev/tour/basics/15
