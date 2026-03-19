# 01 变量与常量（作用域、`_`、`iota`）

> 目标：把“变量/常量的作用域、匿名变量 `_`、`iota` 常量枚举”讲清楚，并能跑通一个最小示例。  
> 适用场景（后台管理 API）：写 handler/service/repo 时，你会频繁遇到“临时变量的可见范围”“忽略不需要的返回值”“用常量表达稳定的角色/状态/错误类别”等。

## 先复习 3 件事（不跳转旧笔记，直接在这里回忆）

1) `var`/`:=` 的直觉
- `var x int`：声明变量，初始化为零值（这里是 `0`）
- `x := 1`：短变量声明（只能在函数体内），= 声明 + 赋值 + 类型推导

2) “未使用即报错”的好处
- 变量/import 没用到会编译失败，逼你把代码保持干净

3) shadowing（遮蔽）
- 在一个 block（大括号）里写 `x := ...`，可能创建一个**新的同名变量**，外层 `x` 不会变

---

## 知识点 1：作用域（scope）+ 块级作用域（block scope）+ `:=` 遮蔽

### B. 一句话定义
作用域就是“变量在哪些代码范围内可见”；Go 里最常见的边界是大括号 `{}`，block 内声明的变量通常只在 block 内有效。

### C. 为什么重要（不做会怎样）
后台管理 API 里你会写大量这种代码：
- handler：解析参数 → 调 service → 返回 JSON
- service：做业务校验 → 调 repo

如果你搞不清作用域/遮蔽，就会出现这种隐蔽 bug：
- 你以为“外层 err/x 已经更新”，其实你在 if/for 的 block 里用 `:=` 新建了一个同名变量，外层还是旧值 → 导致后续逻辑用错变量，排查非常痛苦。

### D. 重难点拆解（2–4 条）
1) **块级作用域**：`if/for/switch` 的 `{}` 里声明的变量，出了 `{}` 就没了。  
2) **`:=` 的规则**：在当前作用域里，`:=` 会“声明新变量”；如果名字已存在且满足“至少一个新变量”的条件，也可能同时“复用 + 新建”，更容易读错。  
3) **推荐 if init 的真实原因**：`if v, err := f(); err != nil { ... }` 是为了把 `v/err` 的生命期限制在 if 内，减少误用；但你得知道它不会影响外层同名变量。  

### E. 业务场景落地
比如你在 handler 里做参数解析：
- 如果解析失败就直接返回 400  
- 如果解析成功才继续调用 service  

这时把 `page,size,err` 控制在最小范围内，能减少“拿错变量继续往下跑”的事故。

### F. 主要代码（关键片段）
```go
x := 1
if true {
	x := 2
	fmt.Printf("inner x=%d\n", x) // Output: inner x=2
}
fmt.Printf("outer x=%d\n", x) // Output: outer x=1
```

### 知识点运用示例
示例 1.1：修复遮蔽，让外层变量真的被更新  
- 要求：把上面片段里的 `x := 2` 改成 `x = 2`  
- 验收：输出里的 `outer x=2`

参考实现（主要代码片段）：
```go
x := 1
if true {
	x = 2
	fmt.Printf("inner x=%d\n", x) // Output: inner x=2
}
fmt.Printf("outer x=%d\n", x) // Output: outer x=2
```

---

## 知识点 2：匿名变量 `_`（blank identifier）

### B. 一句话定义
`_` 是“丢弃桶”：把你不需要的值明确丢掉（常见于多返回值、`range`）。

### C. 为什么重要（不做会怎样）
后台管理 API 常见多返回值：
- `(value, error)`（最常见）
- `(rows, err)` / `(result, err)`（数据库）

如果你不使用 `_`，就会触发“未使用即报错”；如果你滥用 `_`（尤其忽略 error），就会把错误吞掉，线上排查地狱。

### D. 重难点拆解（2–4 条）
1) `_` 不是“随便写写”，它表示“我明确不关心这个值”。  
2) **不要用 `_` 忽略 `error`**（除非你非常确定且能解释原因）；多数时候应该处理/返回/记录。  
3) `range` 里不需要 index 就写 `_`，不需要 value 也可以写 `_`。  

### E. 业务场景落地
比如你调用一个函数只需要 `id` 不需要 `affectedRows`，用 `_` 明确忽略；但对于 `err`，你通常必须处理并统一映射成错误码/HTTP status。

### F. 主要代码（关键片段）
忽略多返回值中的一个值：
```go
q, _ := divMod(10, 3)
fmt.Printf("10/3 quotient=%d\n", q) // Output: 10/3 quotient=3
```

`range` 里忽略 index（只用 value）：
```go
names := []string{"alice", "bob"}
for _, name := range names {
	fmt.Printf("name=%s\n", name) // Output: name=alice（第一次迭代）；Output: name=bob（第二次迭代）
}
```

### 知识点运用示例
示例 2.1：不要丢弃余数，把它也打印出来  
- 要求：把 `q, _ := divMod(10,3)` 改成 `q, rem := ...` 并打印 `rem`  
- 验收：输出里出现 `10%3 remainder=1`

参考实现（主要代码片段）：
```go
q, rem := divMod(10, 3)
fmt.Printf("10/3 quotient=%d\n", q)     // Output: 10/3 quotient=3
fmt.Printf("10%%3 remainder=%d\n", rem) // Output: 10%3 remainder=1
```

---

## 知识点 3：`const` + `iota`（常量计数器）

### B. 一句话定义
`iota` 只在 `const (...)` 块里生效：从 0 开始，每行 +1；每个 const 块都会重置。

### C. 为什么重要（不做会怎样）
后台管理 API 里“稳定枚举/常量”非常多：
- 角色：Admin / Editor / Viewer
- 状态：Enabled / Disabled
- 错误类别：INVALID_PARAM / NOT_FOUND / CONFLICT ...

如果不用 `const`，你就会在代码里散落一堆 magic number/string，后续改名/排查/统一响应会非常痛苦。

### D. 重难点拆解（2–4 条）
1) `iota` **按行递增**：哪怕某行写了显式值，它也会继续按行 +1。  
2) `iota` **每个 const 块从 0 重置**：别把同一组枚举拆散到多个 const 块里。  
3) 枚举常见做法：`iota + 1` 跳过 0（0 有时会留给“未知/未设置”）。  

### E. 业务场景落地
角色枚举是 RBAC 的基础：你会在鉴权/审计日志/管理后台里反复使用；用 `const` 能保证“全局一致、可检索、可重构”。

### F. 主要代码（关键片段）
模式 1：稳定枚举（角色/状态/错误类别）
```go
type Role int

const (
	RoleAdmin Role = iota + 1
	RoleEditor
	RoleViewer
)
```

模式 2：位标志/单位（容量单位是经典例子）
```go
const (
	_ = iota // skip 0
	KB = 1 << (10 * iota)
	MB
	GB
)
```

### 知识点运用示例
示例 3.1：扩展枚举常量，增加 `RoleGuest`  
- 要求：在 Role 的 const 块里加 `RoleGuest`，并把它打印出来  
- 验收：`RoleGuest=4`（连续递增）

参考实现（主要代码片段）：
```go
const (
	RoleAdmin Role = iota + 1
	RoleEditor
	RoleViewer
	RoleGuest
)

fmt.Printf("RoleGuest=%d\n", RoleGuest) // Output: RoleGuest=4
```

---

## 可运行示例（仓库里已有）

运行：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/kp/vars-const-scope-iota
```

典型输出（顺序固定）：
```text
== KP: vars/const/scope/_/iota ==
[1] scope + block scope + := shadowing
inner x=2
outer x=1
[2] blank identifier (_)
10/3 quotient=3
name=alice
name=bob
[3] const + iota
RoleAdmin=1 RoleEditor=2 RoleViewer=3
KB=1024 MB=1048576 GB=1073741824
```

---

## References
- 官方：Go Spec（Constants / Iota / Scopes / Blank identifier）https://go.dev/ref/spec  
- 官方：Effective Go（Constants / iota 常见用法）https://go.dev/doc/effective_go  
- 官方：A Tour of Go（Constants）https://go.dev/tour/basics/15
