# Day 03.1：struct + 方法（值/指针接收者）+ API optional 字段（指针字段）

> 贯穿项目：后台管理 API（RBAC + 登录鉴权 + 列表分页检索 + CRUD + 审计日志 + Docker 部署）  
> 本日目标：把“数据结构 + 行为 + API 可选字段建模”打牢，为后面的 service/repo 分层与 PATCH 接口铺路。

统一运行目录：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
```

---

## 知识点 1：`struct`（数据结构体）+ 零值（zero value）

### B. 一句话定义
`struct` 把一组字段组合成一个类型；字段未赋值就是零值（0/""/false/nil）。

### C. 为什么重要（不做会怎样）
你的后台管理项目里：用户、角色、权限、分页响应、错误响应，几乎都是 `struct`；不理解零值会导致“缺失 vs 有效 0 值”混淆。

### D. 重难点拆解（2–4 条）
1) **零值是语言契约**：不要把零值当“没初始化的垃圾”。  
2) **组合优于继承**：Go 没有 class 继承链，结构体设计更偏“数据 + 组合”。  
3) **输出可读性**：调试常用 `%+v` 看字段名和值。

### E. 业务场景落地
你会有 `User`、`Role`、`Permission`、`ListUsersResponse` 等 struct；零值决定了“缺省行为”。

### F. 代码示例（最小可运行）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day03_01a_struct_zero/main.go`
```go
package main

import "fmt"

type User struct {
	ID    int
	Name  string
	Admin bool
}

func main() {
	fmt.Println("== Day03.1a: struct + zero value ==") // Output: == Day03.1a: struct + zero value ==

	var u User
	fmt.Printf("u=%+v\n", u) // Output: u={ID:0 Name: Admin:false}

	u2 := User{ID: 1, Name: "Alice", Admin: true}
	fmt.Printf("u2=%+v\n", u2) // Output: u2={ID:1 Name:Alice Admin:true}

	u3 := User{Name: "Bob"} // 没赋值的字段会是零值
	fmt.Printf("u3=%+v\n", u3) // Output: u3={ID:0 Name:Bob Admin:false}
}
```

### G. 怎么运行（命令 + 预期现象）
```bash
go run ./cmd/day03_01a_struct_zero
# Output: == Day03.1a: struct + zero value ==
# Output: u={ID:0 Name: Admin:false}
# Output: u2={ID:1 Name:Alice Admin:true}
# Output: u3={ID:0 Name:Bob Admin:false}
```

### H. 练习题（1–3 题）
练习 1：给 `User` 加一个字段 `Age int`，并验证零值输出  
- 验收标准：`Age` 在零值 `u` 中打印为 `0`

### I. 参考答案
参考答案 1：  
把 `Age int` 加进 struct，重新跑 `go run ./cmd/day03_01a_struct_zero`，输出里会出现 `Age:0`（典型输出可略有差异，但零值应为 0）。

---

## 知识点 2：方法（method）+ 值/指针接收者（receiver）取舍

### B. 一句话定义
方法就是“带接收者的函数”；值接收者操作副本，指针接收者操作同一个对象（可修改）。

### C. 为什么重要（不做会怎样）
你的 service/repo 类型几乎都会用方法组织行为；如果你接收者选错，会出现“以为修改了，其实没生效”的隐蔽 bug。

### D. 重难点拆解（2–4 条）
1) **要修改就用指针接收者**：这是最常见且最安全的经验法则。  
2) **避免大对象拷贝**：结构体很大时，值接收者会拷贝成本高。  
3) **Go 会自动取地址/解引用**：调用体验更顺，但你仍要明确语义。

### E. 业务场景落地
例如 `type UserService struct{ repo ... }`，方法 `CreateUser(...)` 基本都用指针接收者；否则你很难在内部维护状态（比如缓存、统计）。

### F. 代码示例（最小可运行）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day03_01b_methods_receivers/main.go`
```go
package main

import "fmt"

type Counter struct {
	N int
}

// 值接收者：操作的是“副本”（copy）
func (c Counter) IncByValue() {
	c.N++
}

// 指针接收者：操作的是“同一个对象”
func (c *Counter) IncByPtr() {
	c.N++
}

func main() {
	fmt.Println("== Day03.1b: methods + value/pointer receiver ==") // Output: == Day03.1b: methods + value/pointer receiver ==

	c := Counter{N: 10}
	c.IncByValue()
	fmt.Printf("after IncByValue: c.N=%d\n", c.N) // Output: after IncByValue: c.N=10

	c.IncByPtr()
	fmt.Printf("after IncByPtr:   c.N=%d\n", c.N) // Output: after IncByPtr:   c.N=11

	// Go 会在需要时自动取地址/解引用，让方法调用更顺滑
	cp := &Counter{N: 20}
	cp.IncByPtr()
	fmt.Printf("cp.N=%d\n", cp.N) // Output: cp.N=21
}
```

### G. 怎么运行（命令 + 预期现象）
```bash
go run ./cmd/day03_01b_methods_receivers
# Output: == Day03.1b: methods + value/pointer receiver ==
# Output: after IncByValue: c.N=10
# Output: after IncByPtr:   c.N=11
# Output: cp.N=21
```

### H. 练习题（1–3 题）
练习 1：把 `IncByValue` 改成返回新值（不改原对象），并在 main 里用返回值更新 `c`  
- 验收标准：你能得到最终 `c.N=11`（注意输出标注规则）

### I. 参考答案
参考答案 1（示例思路）：  
把 `IncByValue()` 改成：
```go
func (c Counter) IncByValue() Counter { c.N++; return c }
```
然后 `c = c.IncByValue()`。

---

## 知识点 3：API optional 字段（PATCH）：用指针区分“缺失 vs 显式零值”

### B. 一句话定义
把字段声明成 `*T`：`nil` 表示“没传”；非 `nil` 表示“传了（即使值是零值）”。

### C. 为什么重要（不做会怎样）
后台管理里“编辑用户/部分更新”非常常见：如果你用基础类型，你分不清“没传 admin”与“传了 admin=false”，会把数据改错。

### D. 重难点拆解（2–4 条）
1) **`omitempty` 的语义**：`nil` 指针会被省略；非 nil 即使是 false/0 也会出现。  
2) **只更新非 nil 字段**：这是 PATCH 的核心逻辑。  
3) **不要把 optional 滥用到响应**：响应一般用非指针，让调用方更省心（除非你明确要表达“未知/缺失”）。

### E. 业务场景落地
`PATCH /api/v1/users/{id}`：前端可能只传 `{ "admin": false }`，你必须把它当成“显式设置”，而不是“缺省”。

### F. 代码示例（最小可运行）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day03_01c_json_optional_fields/main.go`
```go
package main

import (
	"encoding/json"
	"fmt"
)

// 在 HTTP API 里，为了区分“字段缺失” vs “字段显式给了零值”，常用指针字段表示 optional。
type PatchUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Admin *bool   `json:"admin,omitempty"`
}

func main() {
	fmt.Println("== Day03.1c: JSON optional fields (pointer) ==") // Output: == Day03.1c: JSON optional fields (pointer) ==

	rawMissing := []byte(`{}`)
	var r1 PatchUserRequest
	_ = json.Unmarshal(rawMissing, &r1)
	fmt.Printf("missing: nameNil=%v adminNil=%v\n", r1.Name == nil, r1.Admin == nil) // Output: missing: nameNil=true adminNil=true

	rawZero := []byte(`{"admin": false}`)
	var r2 PatchUserRequest
	_ = json.Unmarshal(rawZero, &r2)
	fmt.Printf("explicit false: adminNil=%v adminVal=%v\n", r2.Admin == nil, *r2.Admin) // Output: explicit false: adminNil=false adminVal=false

	name := "Alice"
	admin := true
	out, _ := json.Marshal(PatchUserRequest{Name: &name, Admin: &admin})
	fmt.Printf("marshal: %s\n", string(out)) // Output: marshal: {"name":"Alice","admin":true}
}
```

### G. 怎么运行（命令 + 预期现象）
```bash
go run ./cmd/day03_01c_json_optional_fields
# Output: == Day03.1c: JSON optional fields (pointer) ==
# Output: missing: nameNil=true adminNil=true
# Output: explicit false: adminNil=false adminVal=false
# Output: marshal: {"name":"Alice","admin":true}
```

### H. 练习题（1–3 题）
练习 1：写一个 `ApplyPatch`：只更新非 nil 字段（纯内存，不接 DB）  
- 验收标准：
  - 传 `{}` 不改变 user
  - 传 `{"admin": false}` 会把 user.Admin 设置为 false

### I. 参考答案
参考答案 1（可运行做法）：  
你可以先在 `day03_01c_json_optional_fields/main.go` 里加一个 `User` struct 和 `ApplyPatch(user, req)` 函数，最后 `fmt.Printf` 打印更新前后（打印要写典型输出注释）。

---

## References
- 官方：A Tour of Go（Structs）https://go.dev/tour/moretypes/2
- 官方：A Tour of Go（Methods）https://go.dev/tour/methods/1
- 官方：Effective Go https://go.dev/doc/effective_go
