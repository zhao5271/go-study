# Day 03.1：struct 与方法（值/指针接收者）+ API 里的 optional 字段

> 目标：把 TS/Node 的“对象/类/接口”心智迁移到 Go 的 `struct` 与方法；并引入 API 常用的“可选字段”建模方式（指针字段）。

## 1) 知识讲解：概念 → 为什么（设计动机/取舍）

### 1.1 `struct` 是 Go 的“数据结构体”（偏数据，不偏继承）
- `struct` 用来把一组字段组合成一个类型（更像 TS 的 object type + 数据载体）。
- Go 没有 class 的继承链；更鼓励“组合（composition）”而非“继承（inheritance）”。
- 字段没赋值会是 **零值（zero value）**：`int` 是 `0`、`string` 是 `""`、`bool` 是 `false`、指针是 `nil`。

### 1.2 方法（method）= “函数 + 接收者（receiver）”
- Go 的方法本质是“带接收者的函数”，用于给类型绑定行为。
- 关键取舍：**值接收者 vs 指针接收者**
  - 值接收者：拿到的是副本（copy），适合不修改接收者、或者类型很小且逻辑明确的场景。
  - 指针接收者：能修改接收者、避免大对象拷贝，通常更常见。

### 1.3 在 API 里，optional 字段常用“指针字段”区分“缺失 vs 显式零值”
- TS 里你常用 `undefined` 表示缺失；Go 的基础类型没有 `undefined`，缺省就是零值。
- 对 PATCH/部分更新场景，常见需求是区分：
  - “客户端没传 admin 字段” vs “客户端明确传了 admin=false”
- 解决方式：把字段改成指针：`*bool`、`*int`、`*string`，用 `nil` 表示缺失。

## 2) 示例驱动：每个知识点后立刻可运行代码（全文）

### 2.1 `struct` + 零值 + 复合字面量

运行：
- `cd go-learning && go run ./cmd/day03_01a_struct_zero`

代码（全文）：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day03_01a_struct_zero/main.go`
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

### 2.2 方法 + 值/指针接收者（是否能修改接收者）

运行：
- `cd go-learning && go run ./cmd/day03_01b_methods_receivers`

代码（全文）：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day03_01b_methods_receivers/main.go`
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

### 2.3 API optional 字段（PATCH 思维）：用指针区分“缺失 vs 显式零值”

运行：
- `cd go-learning && go run ./cmd/day03_01c_json_optional_fields`

代码（全文）：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day03_01c_json_optional_fields/main.go`
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

## 3) 常见坑（结合 TS/Node 习惯对照）

1) **以为“值接收者也能改原对象”**
- TS 里对象默认引用语义；Go 里 `struct` 值是“按值拷贝”。
- 经验法则：只要方法需要修改接收者，直接用指针接收者 `func (x *T) ...`。

2) **把“字段缺失”误当成“字段为零值”**
- TS 的 `undefined` 很常见；Go 基础类型没这个概念。
- 做 PATCH/可选字段时，用指针字段（`*T`）或 `sql.Null*`/自定义 Optional 类型表达“是否出现”。

3) **在切片里存 `struct`，在循环里修改却没生效**
- `for _, v := range sliceOfStruct { v.Field = ... }` 修改的是副本。
- 解决：用索引遍历 `for i := range s { s[i].Field = ... }` 或存指针 `[]*T`。

## 4) 工程用法/最佳实践：真实 API 项目怎么落地

- 分层建模（对齐你 Node/Java 经验）：
  - `transport/http`：请求/响应 DTO（带 `json` tag），与路由/校验强绑定
  - `service`：业务输入/输出结构（更贴业务语义，尽量不暴露 transport 细节）
  - `repo`：DB 实体与查询结果（字段更贴数据库）
- PATCH 语义建议（配合 `api-design-principles` 的“清晰一致”原则）：
  - 如果是“部分更新”，请求体里用指针字段区分“缺失”与“显式零值”
  - 响应体一般用非指针字段（客户端更省心），除非你确实要表达“未知/缺失”

## 5) 练习策略：可直接作为“运用示例”（含完整参考实现与讲解）

练习目标：做一个“用户 PATCH 更新”的纯内存版本（先不接 MySQL），把 optional 字段逻辑跑通。

要求（你自己写时建议对照这份参考实现）：
- 入参：`PatchUserRequest`（指针字段）
- 逻辑：只更新非 nil 的字段
- 返回：更新后的用户

下一小步我会带你把它抽成 `service` 包（不引入 `go test`，只用 `go run` 验证）。

## References
- [A Tour of Go: Structs](https://go.dev/tour/moretypes/2)（官方：struct 基础与字面量）
- [A Tour of Go: Methods](https://go.dev/tour/methods/1)（官方：方法与接收者）
- [Effective Go](https://go.dev/doc/effective_go)（官方：Go 风格与工程化建议）
