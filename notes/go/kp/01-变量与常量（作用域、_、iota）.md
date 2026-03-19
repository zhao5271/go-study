# 01 变量与常量（作用域、`_`、`iota`）

> 更新于：2026-03-19  
> 目标：把“变量/常量的作用域、匿名变量 `_`、`iota` 常量枚举”讲清楚，并能在后台管理 API（RBAC/错误码/分页）里正确落地。

## TL;DR（可放入 progress/context-pack）

1) Go 是**词法作用域**：变量只在声明位置所属的 `{}`（block）里可见。  
2) `if/for/switch` 的 **init 语句**里声明的变量，作用域只在该语句块内（包含 `else` / `case`）。  
3) `:=` 可能触发 **shadowing（遮蔽）**：内层同名变量会让你“以为改了外层”，但其实没改。  
4) `_` 是“丢弃桶”（blank identifier）：用于显式丢弃不需要的值；**不要用 `_` 吞 error**。  
5) `const` 是编译期常量；`iota` 只在 `const (...)` 里生效：从 0 开始，**按行递增**，每个 const 块会**重置**。  
6) 工程上：RBAC 角色/状态/错误码用 `type + const(iota)`；权限用位标志 `1 << iota`。  
7) 为了可维护：枚举常量尽量**带类型**（`type Role int`），避免 magic number/string 漫天飞。  

## 关键词
- scope / block scope / lexical scope（作用域/块级/词法）  
- short variable declaration `:=`（短变量声明）  
- shadowing（遮蔽）  
- blank identifier `_`（空标识符）  
- `const` / `iota`（常量 / 常量计数器）  
- enum（枚举）、bit flag（位标志）  

## 关联复习（不跳转旧笔记）

- `var` 声明的变量有零值；`:=` 是“声明 + 赋值 + 类型推导”（只能在函数体内）。  
- Go 编译期会报“未使用”的变量/import，这是为了逼你保持代码干净（也减少“写了但忘了用”的 bug）。  
- 函数常见多返回值 `(value, error)`：**error 是返回值的一部分，不是异常**。  

---

## 知识点 1：作用域（scope）+ `:=` 遮蔽（shadowing）

### 一句话定义
作用域就是“一个名字在代码里**能被引用**的范围”；在 Go 里最常见的边界是大括号 `{}`（块级作用域）。

### 为什么重要（工程视角）
后台管理 API 里，handler/service/repo 都充满“短生命周期临时变量”（分页参数、校验结果、DB 查询返回值）。  
作用域/遮蔽没搞清楚，会出现最难排查的一类 bug：**你以为外层变量被更新了，其实你在内层用 `:=` 新建了同名变量**。

### 关键坑与取舍（2–4 条）
1) **init 语句的作用域更“窄”**：`if x, err := ...; err != nil { ... }` 里的 `x/err` 只活在该 if（包含 else）里，出了 if 就没了。  
2) **`:=` 的“至少一个新变量”规则**：左边如果有“至少一个新名字”，这条语句就会在当前作用域里重新声明所有左侧名字（更容易不小心遮蔽）。  
3) **需要更新外层变量时用 `=`**：在内层 block 里写 `x = ...` 才是更新外层；写 `x := ...` 通常是在创建新变量。  
4) 取舍：`if init` 能缩小变量生命周期，减少误用；但如果你本来就需要复用外层变量（比如外层 `err`），就要避免在内层再 `:=` 同名。  

### 业务场景落地（后台管理 API）
分页解析是高频代码：`page/size` 从 query（string）来 → 转换/校验 → 计算 offset。  
如果你把默认值写在外层，再在 if 里用 `:=`，很容易“解析成功但外层没更新”，导致永远用默认分页。

### 主要代码（最小示例）

遮蔽导致“看起来更新了，其实没更新”：
```go
pageSize := 10
raw := "20"
if v, err := strconv.Atoi(raw); err == nil {
	pageSize := v // 坑：遮蔽外层 pageSize（只在 if 内有效）
	fmt.Printf("inner pageSize=%d\n", pageSize) // Output: inner pageSize=20
}
fmt.Printf("outer pageSize=%d\n", pageSize) // Output: outer pageSize=10
```

正确写法（更新外层变量）：
```go
pageSize := 10
raw := "20"
if v, err := strconv.Atoi(raw); err == nil {
	pageSize = v
	fmt.Printf("inner pageSize=%d\n", pageSize) // Output: inner pageSize=20
}
fmt.Printf("outer pageSize=%d\n", pageSize) // Output: outer pageSize=20
```

### 知识点运用示例

练习 1.1：修复遮蔽 bug  
- 要求：把“遮蔽示例”里 `pageSize := v` 改成正确写法，让外层 pageSize 也更新。  
- 验收：`outer pageSize=20`。  

参考答案（可直接替换）：
```go
pageSize := 10
raw := "20"
if v, err := strconv.Atoi(raw); err == nil {
	pageSize = v
	fmt.Printf("inner pageSize=%d\n", pageSize) // Output: inner pageSize=20
}
fmt.Printf("outer pageSize=%d\n", pageSize) // Output: outer pageSize=20
```

---

## 知识点 2：空标识符 `_`（blank identifier）

### 一句话定义
`_` 是一个“只能写、不能读”的名字：用来显式丢弃你不需要的值。

### 为什么重要（工程视角）
Go 的多返回值在 API/DB 场景极其常见；`_` 让你“明确表示不关心某个返回值”，同时避免未使用编译错误。  
但如果你用 `_` 吞掉 `error`，你就等于把“失败路径”删掉了：线上出问题时定位会非常困难。

### 关键坑与取舍（2–4 条）
1) `_` 的语义是“我明确不关心”，不是“先写着凑编译”。`_ = x` 只能临时救火，最终要删掉。  
2) **不要用 `_` 忽略 error**：error 是业务流程的一部分（例如 DB 连接失败、参数非法），忽略会让程序在错误状态继续运行。  
3) `range` 里常用 `_` 忽略 index：`for _, v := range xs { ... }`。  
4) 典型工程用法（非常推荐记住）：
   - **类型断言只关心 ok**：`_, ok := v.(T)`  
   - **编译期接口校验**：`var _ Interface = (*Impl)(nil)`（只用于编译期保证“实现了接口”，不会在运行期创建对象）  

### 业务场景落地（后台管理 API）
- RBAC：你可能只需要判断“是不是管理员/有没有某个权限”，不需要把整对象都拿出来。  
- 工程化：你写 repo/service 分层时，用“编译期接口校验”可以避免重构后悄悄漏实现导致的运行期 bug。  

### 主要代码（最小示例）

忽略 `range` 的 index（只用 value）：
```go
names := []string{"alice", "bob"}
for _, name := range names {
	fmt.Printf("name=%s\n", name) // Output: name=alice（第一次迭代）；Output: name=bob（第二次迭代）
}
```

类型断言只关心是否是某种类型（不关心值）：
```go
var v any = "admin"
_, ok := v.(string)
fmt.Printf("isString=%v\n", ok) // Output: isString=true
```

编译期接口校验（常用于 repo/service 装配时防止漏实现）：
```go
type UserRepo interface{ Ping() error }
type MySQLUserRepo struct{}
func (*MySQLUserRepo) Ping() error { return nil }

var _ UserRepo = (*MySQLUserRepo)(nil) // Output: （无输出；编译期校验，不会运行）
```

### 知识点运用示例

练习 2.1：把“只取商、不取余数”的代码改成同时拿到余数  
- 要求：把 `q, _ := divMod(10,3)` 改成 `q, rem := ...` 并打印 rem。  
- 验收：输出出现 `10%3 remainder=1`。  

参考答案：
```go
q, rem := divMod(10, 3)
fmt.Printf("10/3 quotient=%d\n", q)     // Output: 10/3 quotient=3
fmt.Printf("10%%3 remainder=%d\n", rem) // Output: 10%3 remainder=1
```

练习 2.2：禁止吞 error  
- 要求：找一个返回 `(int, error)` 的解析函数（例如 `strconv.Atoi`），写出“错误做法（用 _ 忽略 err）”与“正确做法（处理 err）”。  
- 验收：正确做法在 err!=nil 时不会继续使用返回值（你可以打印/返回/写成 400 错误码）。  

参考答案（示例：打印并提前返回）：
```go
raw := "not-a-number"
n, err := strconv.Atoi(raw)
if err != nil {
	fmt.Printf("parse failed: %v\n", err) // Output: parse failed: （输出可能变化/不固定：不同 Go 版本/平台错误文本可能略有差异）
	return
}
fmt.Printf("n=%d\n", n) // Output: （不会执行到这里）
```

---

## 知识点 3：`const` + `iota`（常量枚举/位标志）

### 一句话定义
`const` 声明编译期常量；`iota` 是常量计数器，只在 `const (...)` 块里有效：从 0 开始，按行递增，每个 const 块重置。

### 为什么重要（工程视角）
后台管理 API 里“稳定枚举”非常多：角色/状态/错误码/权限位。  
用 `const + iota` 能把这些“业务约定”固化成可检索、可重构的代码，不再依赖散落各处的 magic number/string。

### 关键坑与取舍（2–4 条）
1) `iota` **按行递增**：哪怕你某行写了显式值，下一行的 `iota` 仍然会按行数递增。  
2) `iota` **每个 const 块重置**：同一组枚举尽量放在一个 const 块，别拆散。  
3) **表达式复用**：const 块里后续行省略表达式，会复用上一行的表达式与类型（常用于单位/位移）。  
4) 取舍：建议用“带类型的枚举”（`type Role int`），好处是更不容易把不同枚举混用；代价是必要时需要显式转换（但这通常是好事）。  

### 业务场景落地（后台管理 API）
- RBAC 角色：`RoleAdmin/RoleEditor/...`（建议 `iota + 1` 跳过 0，把 0 留给 Unknown）。  
- 权限位：`READ/WRITE/EXPORT` 最适合位标志（可以组合），在 DB 里也容易存储（一个整数）。  
- 错误码枚举：统一错误码表时，常用 const 枚举让前后端对齐。  

### 主要代码（最小示例）

角色枚举（跳过 0）：
```go
type Role int

const (
	RoleAdmin Role = iota + 1
	RoleEditor
	RoleViewer
)

fmt.Printf("RoleAdmin=%d RoleEditor=%d RoleViewer=%d\n", RoleAdmin, RoleEditor, RoleViewer) // Output: RoleAdmin=1 RoleEditor=2 RoleViewer=3
```

权限位标志（可组合）：
```go
type Perm uint64

const (
	PermRead Perm = 1 << iota
	PermWrite
	PermExport
)

mask := PermRead | PermExport
fmt.Printf("hasRead=%v\n", (mask&PermRead) != 0)     // Output: hasRead=true
fmt.Printf("hasWrite=%v\n", (mask&PermWrite) != 0)   // Output: hasWrite=false
fmt.Printf("hasExport=%v\n", (mask&PermExport) != 0) // Output: hasExport=true
```

单位（表达式复用 + skip 0）：
```go
const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

fmt.Printf("KB=%d MB=%d GB=%d\n", KB, MB, GB) // Output: KB=1024 MB=1048576 GB=1073741824
```

### 知识点运用示例

练习 3.1：扩展角色枚举  
- 要求：在 Role 的 const 块里加 `RoleGuest`，并打印它。  
- 验收：`RoleGuest=4`（连续递增）。  

参考答案：
```go
const (
	RoleAdmin Role = iota + 1
	RoleEditor
	RoleViewer
	RoleGuest
)
fmt.Printf("RoleGuest=%d\n", RoleGuest) // Output: RoleGuest=4
```

练习 3.2：写一个 `hasPerm(mask, p)` 工具函数  
- 要求：实现 `hasPerm(mask, p Perm) bool`，用于判断某个权限位是否存在。  
- 验收：对 `mask := PermRead|PermExport`，`hasPerm(mask, PermWrite)` 返回 false。  

参考答案：
```go
func hasPerm(mask, p Perm) bool { return (mask&p) != 0 }

mask := PermRead | PermExport
fmt.Printf("hasWrite=%v\n", hasPerm(mask, PermWrite)) // Output: hasWrite=false
```

---

## 可运行示例（仓库里已有）

运行：
```bash
cd go-learning
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

## References
- 官方：Go Spec（Scopes / Declarations / Constants / Iota / Blank identifier）https://go.dev/ref/spec（官方）  
- 官方：Effective Go（Constants / iota 常见用法）https://go.dev/doc/effective_go（官方）  
- 官方：A Tour of Go（Constants）https://go.dev/tour/basics/15（官方）
