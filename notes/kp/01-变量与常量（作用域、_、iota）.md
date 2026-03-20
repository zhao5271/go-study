---
type: kp
domain: go
topic: vars-const-scope-iota
topic_zh: 变量与常量
stage: foundation
status: evergreen
review_cycle: weekly
source:
  - official
tags:
  - go
  - kp
  - foundation
created: 2026-03-19
updated: 2026-03-20
---

# 01 变量与常量（作用域、`_`、`iota`）

> 更新于：2026-03-20  
> 目标：把“作用域、blank identifier、常量枚举”整理成可长期复习的知识库笔记，并能落到后台管理 API 的参数解析、错误处理、RBAC/状态码建模里。

## TL;DR（可放入 progress/context-pack）
- Go 是词法作用域；变量在最内层 block 内生效，内层同名声明会遮蔽外层变量。
- `:=` 既可能声明新变量，也可能“部分重声明”；写错最容易制造 shadowing。
- `_` 是 blank identifier，可显式丢弃不需要的值，但不要吞掉真正需要处理的 `error`。
- `const` 是编译期常量；`iota` 只在 `const (...)` 块里按行递增，每个块从 0 重新开始。
- 工程里角色、状态、错误类别适合用 `type + const + iota`；权限位适合 `1 << iota`。
- 这类基础语义一旦写错，最常见后果不是“编译不过”，而是“逻辑看着对，运行结果悄悄偏了”。

## 关键词
- scope / lexical scope / block scope
- short variable declaration / shadowing
- blank identifier
- const / iota
- enum / bit flag

## 知识点 1：作用域与 `:=` 遮蔽（shadowing）

### 一句话定义
作用域就是“名字在哪一段代码里有效”；Go 使用词法作用域，最常见边界是大括号 `{}`。

### 为什么重要
后台管理 API 里会频繁出现 `page`、`size`、`err`、`user` 这类临时变量；一旦在 `if/for/switch` 里误用 `:=` 新建同名变量，你以为改的是外层值，实际改的是内层副本。

### 重难点拆解
- 函数内部声明的变量，作用域从声明结束后开始，到最内层 block 结束。
- `:=` 只能在函数体内使用；同一作用域里必须至少有一个新变量。
- 内层 block 可以重名声明外层变量；内层名字生效期间，外层同名变量被遮蔽。
- `if err := ...; err != nil {}` 是合理缩小作用域；但如果你后面还要复用外层 `err`，就要警惕误遮蔽。

### 业务场景落地
做用户列表接口时，你经常会先解析 query，再调 service：
- 适合缩小作用域：`if page, err := parsePage(...); err != nil { ... }`
- 不适合遮蔽：外层已经有 `err`，但你在内层再写 `:=`，导致日志里看到的不是你以为的那个错误值

### 代码示例
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

### 怎么运行
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/kp/vars-const-scope-iota
# 预期会看到：
# == KP: vars/const/scope/_/iota ==
# [1] scope + block scope + := shadowing
# inner x=2
# outer x=1
```

### 练习题 1
补全下面代码，让最终打印结果是 `outer x=2`，不要新增第二个 `x`。

**验收标准**
- 程序可运行
- 最终输出 `outer x=2` // Output: outer x=2

**参考答案**
```go
package main

import "fmt"

func main() {
	x := 1
	if true {
		x = 2
	}
	fmt.Printf("outer x=%d\n", x) // Output: outer x=2
}
```

## 知识点 2：blank identifier（`_`）

### 一句话定义
`_` 是匿名占位符：它能接收值，但不会创建可用绑定。

### 为什么重要
Go 很多函数会返回多个值；你经常只关心其中一部分。`_` 可以让“我明确不需要这个值”变得可读，但一旦拿它吞掉关键 `error`，就会把真实问题藏起来。

### 重难点拆解
- `_` 常见于多返回值、`for range`、只为副作用导入包等场景。
- `_` 可以在赋值和短变量声明中出现，但它本身不引入新变量。
- `_` 丢掉普通返回值通常没问题；丢掉 `error` 通常是坏味道。
- `var _ Interface = (*Impl)(nil)` 是编译期接口校验，不是运行期创建对象。

### 业务场景落地
后台管理 API 里：
- `id, _ := strconv.Atoi(s)` 是危险写法，因为非法输入会被你静默忽略
- `for _, user := range users` 很常见，因为你只需要值，不需要索引
- `_ "net/http/pprof"` 这种导入用于副作用，适合调试场景

### 代码示例
```go
package main

import "fmt"

func divMod(a, b int) (int, int) {
	return a / b, a % b
}

func main() {
	q, _ := divMod(10, 3)
	fmt.Printf("quotient=%d\n", q) // Output: quotient=3

	names := []string{"alice", "bob"}
	for _, name := range names {
		fmt.Printf("name=%s\n", name) // Output: name=alice（第一次迭代）；name=bob（第二次迭代）
	}
}
```

### 怎么运行
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/kp/vars-const-scope-iota
# 预期会看到：
# [2] blank identifier (_)
# 10/3 quotient=3
# name=alice
# name=bob
```

### 练习题 2
下面两种写法里，哪一种更适合线上代码？为什么？

```go
id, _ := parseUserID(raw)
```

```go
id, err := parseUserID(raw)
if err != nil {
	return err
}
```

**验收标准**
- 能说明为什么第一种写法风险更高
- 能指出 `error` 被忽略会让什么问题更难排查

**参考答案**
- 第二种更适合线上代码。
- 第一种把 `error` 丢给 `_`，会让非法参数、数据脏值、边界输入都被静默吞掉，最终可能演变成“查错对象”“查空数据”这类更难排查的问题。

## 知识点 3：`const` 与 `iota`

### 一句话定义
`const` 表示编译期常量；`iota` 是 `const (...)` 块里的行号计数器，从 0 开始递增。

### 为什么重要
在后台管理 API 里，角色、状态、错误类别都需要“稳定、可读、少魔法值”的定义方式；`iota` 正是 Go 里最常用的枚举常量构造器。

### 重难点拆解
- `iota` 只在当前 `const (...)` 块内有效；新块会重置为 0。
- `iota` 是按行递增，不是按“是否用了它”递增；某一行即使没用到，计数仍会前进。
- 如果 `0` 代表“未设置/非法”，常见写法是从 `iota + 1` 开始。
- 位标志常见写法是 `1 << iota`；适合组合权限，但不适合表达互斥状态。

### 业务场景落地
- `RoleAdmin / RoleEditor / RoleViewer` 适合用枚举常量
- `PermRead / PermWrite / PermExport` 适合用位标志
- 错误类别（如参数错误、鉴权错误、资源不存在）也适合用带类型的常量集中定义

### 代码示例
```go
package main

import "fmt"

type Role int

const (
	RoleAdmin Role = iota + 1
	RoleEditor
	RoleViewer
)

const (
	_ = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

func main() {
	fmt.Printf("roles=%d,%d,%d\n", RoleAdmin, RoleEditor, RoleViewer) // Output: roles=1,2,3
	fmt.Printf("KB=%d MB=%d GB=%d\n", KB, MB, GB)                    // Output: KB=1024 MB=1048576 GB=1073741824
}
```

### 怎么运行
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/kp/vars-const-scope-iota
# 预期会看到：
# [3] const + iota
# RoleAdmin=1 RoleEditor=2 RoleViewer=3
# KB=1024 MB=1048576 GB=1073741824
```

### 练习题 3
定义一个 `AuditStatus` 类型，要求：
- `0` 保留为“未设置”
- `1` 表示 `Pending`
- `2` 表示 `Success`
- `3` 表示 `Failed`

**验收标准**
- 使用带类型的常量
- 不出现 magic number

**参考答案**
```go
package main

type AuditStatus int

const (
	AuditStatusPending AuditStatus = iota + 1
	AuditStatusSuccess
	AuditStatusFailed
)
```

## 关联复习
- `02-基础类型（转换、格式化、表达式）` 里会继续展开显式类型转换、`string(65)` 这类常见误解。
- `03-字符串基本操作（转义、格式化、Builder、比较、常用方法）` 里会继续展开格式化输出和字符串处理。

## References
- [Go Spec](https://go.dev/ref/spec) - 作用域、blank identifier、short variable declaration、`iota`、常量与字符串的官方语言定义（官方）
- [Effective Go](https://go.dev/doc/effective_go) - `_` 的典型用法、不要忽略 `error`、`iota` 的常见模式（官方）
- NotebookLM 查询失败（原因：浏览器 profile 被占用，`ProcessSingletonLock` 冲突），本节结论以官方资料补齐。
