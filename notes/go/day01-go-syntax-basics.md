# Day 01：语法地基（可运行示例 + 可复习笔记）

> 贯穿项目：后台管理 API（RBAC + 登录鉴权 + 列表分页检索 + CRUD + 审计日志 + Docker 部署）  
> 本日目标：你先把 Go 的“工程入口/包机制/零值/错误处理/集合类型”跑通，后面接 HTTP/MySQL 才不会踩基础坑。

统一运行目录：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
```

---

## 知识点 1：`package` / `import` / 导出规则（首字母大小写）+ 工程入口（cmd）

### B. 一句话定义
Go 用 **package** 组织代码；标识符 **首字母大写**表示包外可见；`cmd/` 里放可运行入口。

### C. 为什么重要（不做会怎样）
后台管理 API 很快会出现：handler/service/repo/config 等多模块；如果你不理解 package 边界和可见性，工程会变成“随便互相引用”，后期重构成本爆炸。

### D. 重难点拆解（2–4 条）
1) **未使用 import/变量会编译失败**：Go 用编译器强制你保持依赖干净。  
2) **导出规则是约定而不是关键字**：只看首字母大小写。  
3) **入口与库分离**：`cmd/` 只负责启动/演示，复用逻辑放 `internal/`（后面会用到）。

### E. 业务场景落地
你后面会写 `cmd/server` 启动 HTTP 服务，同时 `internal/api`/`internal/service` 等包作为可复用逻辑；这个分层思维从 Day01 就开始。

### F. 代码示例（最小可运行）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day01_01_packages_import/main.go`
```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("== Day01.1: package/import + exported names ==") // Output: == Day01.1: package/import + exported names ==

	s := "go fullstack"
	fmt.Printf("ToUpper(%q)=%q\n", s, strings.ToUpper(s)) // Output: ToUpper("go fullstack")="GO FULLSTACK"

	// In Go: identifiers starting with Uppercase are exported (public).
	// identifiers starting with lowercase are unexported (package-private).
}
```

### G. 怎么运行（命令 + 预期现象）
```bash
go run ./cmd/day01_01_packages_import
# Output: == Day01.1: package/import + exported names ==
# Output: ToUpper("go fullstack")="GO FULLSTACK"
```

### H. 练习题（1–3 题）
练习 1：故意加一个没用到的 import（比如 `math`），观察编译错误  
- 验收标准：你能看到编译器明确提示 “imported and not used”

### I. 参考答案（紧跟每道练习题）
参考答案 1（可运行做法）：  
在 `main.go` 里加 `import "math"` 且不使用，然后运行：
```bash
go run ./cmd/day01_01_packages_import
# Output: ... imported and not used: "math" (输出可能变化/不固定：报错细节随 Go 版本变化)
```

---

## 知识点 2：变量声明 + 零值（Zero Value）+ `:=`（短变量声明）

### B. 一句话定义
Go 里“声明即初始化为零值”，函数体内常用 `:=` 做“声明+赋值+类型推导”。

### C. 为什么重要（不做会怎样）
你后面写 handler/service 时会大量用到“零值 + 显式错误处理”；如果你把零值当成“未定义”，就会把 bug 带到线上（例如把 `0` 当成“未传”）。

### D. 重难点拆解（2–4 条）
1) **没有 `undefined`**：只有“零值/非零值”，需要你明确“零值是否是合法业务值”。  
2) **`:=` 只能在函数体内**：包级别必须用 `var`。  
3) **指针零值是 `nil`**：这是你区分“缺失/未初始化”的关键手段（后面 JSON PATCH 会用）。

### E. 业务场景落地
例如“分页 page/size”里 `0` 通常不是合法值；你要显式校验，而不是指望“未传就是 0 然后自动变成默认值”。

### F. 代码示例（最小可运行）
文件 1：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day01_02_vars_zero/main.go`
```go
package main

import "fmt"

func main() {
	fmt.Println("== Day01.2: var + zero values ==") // Output: == Day01.2: var + zero values ==

	var count int
	var name string
	var ok bool
	fmt.Printf("count=%d name=%q ok=%v\n", count, name, ok) // Output: count=0 name="" ok=false

	var p *int
	fmt.Printf("p==nil? %v\n", p == nil) // Output: p==nil? true
}
```

文件 2：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day01_03_short_decl/main.go`
```go
package main

import "fmt"

func main() {
	fmt.Println("== Day01.3: := short declaration ==") // Output: == Day01.3: := short declaration ==

	x := 42
	fmt.Printf("x=%d type=%T\n", x, x) // Output: x=42 type=int

	// := is function-scope only; package-level must use var.
}
```

### G. 怎么运行（命令 + 预期现象）
```bash
go run ./cmd/day01_02_vars_zero
# Output: == Day01.2: var + zero values ==
# Output: count=0 name="" ok=false
# Output: p==nil? true

go run ./cmd/day01_03_short_decl
# Output: == Day01.3: := short declaration ==
# Output: x=42 type=int
```

### H. 练习题（1–3 题）
练习 1：把 `var p *int` 改成 `p := new(int)` 并打印 `*p`  
- 验收标准：你能看到 `*p` 的典型输出是 `0`（因为 int 的零值是 0）

### I. 参考答案
参考答案 1（可运行做法）：  
在 `day01_02_vars_zero/main.go` 里加：
```go
p := new(int)
fmt.Printf("*p=%d\n", *p) // Output: *p=0
```
然后运行 `go run ./cmd/day01_02_vars_zero`。

---

## 知识点 3：显式类型转换 + `const` + `(value, error)`（工程里最常见的失败建模）

### B. 一句话定义
Go 不做隐式数值转换；函数失败通常用 `(value, error)` 表达，调用点显式检查 `err`。

### C. 为什么重要（不做会怎样）
后台管理 API 的每一层（handler/service/repo）都会返回错误；如果你用“随便 panic/忽略 err”，系统会变得不可控、不可定位、不可交付。

### D. 重难点拆解（2–4 条）
1) **显式转换**：`int`/`int64` 不会自动转换，避免精度/溢出悄悄发生。  
2) **`const` 是编译期常量**：常用来表达稳定配置/枚举值（后面错误码会用）。  
3) **错误是值**：`ErrXxx`（sentinel error）+ `errors.Is` 是最小可复用模式。

### E. 业务场景落地
例如“创建用户”如果 email 重复，你会返回一个可匹配的错误语义（后面映射成 HTTP 409），而不是靠字符串判断。

### F. 代码示例（最小可运行）
文件 1：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day01_04_conversions_const/main.go`
```go
package main

import "fmt"

func main() {
	fmt.Println("== Day01.4: explicit conversions + const ==") // Output: == Day01.4: explicit conversions + const ==

	x := 42
	var big int64 = 1
	sum := int64(x) + big
	fmt.Printf("int64(x)+big=%d type=%T\n", sum, sum) // Output: int64(x)+big=43 type=int64

	const pi = 3.14159
	fmt.Printf("pi=%.2f type=%T\n", pi, pi) // Output: pi=3.14 type=float64
}
```

文件 2：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day01_05_functions_error/main.go`
```go
package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println("== Day01.5: functions + (value, error) ==") // Output: == Day01.5: functions + (value, error) ==

	v, err := divide(10, 2)
	fmt.Printf("divide(10,2) => v=%d err=%v\n", v, err) // Output: divide(10,2) => v=5 err=<nil>

	_, err = divide(10, 0)
	fmt.Printf("divide(10,0) => err=%v\n", err) // Output: divide(10,0) => err=divide by zero

	fmt.Printf("errors.Is(err, ErrDivideByZero)=%v\n", errors.Is(err, ErrDivideByZero)) // Output: errors.Is(err, ErrDivideByZero)=true
}

var ErrDivideByZero = errors.New("divide by zero")

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}
```

### G. 怎么运行（命令 + 预期现象）
```bash
go run ./cmd/day01_04_conversions_const
# Output: == Day01.4: explicit conversions + const ==
# Output: int64(x)+big=43 type=int64
# Output: pi=3.14 type=float64

go run ./cmd/day01_05_functions_error
# Output: == Day01.5: functions + (value, error) ==
# Output: divide(10,2) => v=5 err=<nil>
# Output: divide(10,0) => err=divide by zero
# Output: errors.Is(err, ErrDivideByZero)=true
```

### H. 练习题（1–3 题）
练习 1：给 `divide` 增加一个分支：`a < 0 || b < 0` 返回新错误 `ErrNegativeNotAllowed`  
- 验收标准：你能用 `errors.Is(err, ErrNegativeNotAllowed)` 匹配到它

练习 2：跑一遍控制流与集合例子，写下“map 顺序为什么不稳定”  
- 验收标准：你能观察到 `day01_07_slices_maps` 的遍历输出行有“输出可能变化/不固定”标注

### I. 参考答案
参考答案 1（可运行做法）：  
按练习描述新增：
```go
var ErrNegativeNotAllowed = errors.New("negative not allowed")
```
并在 `divide` 里判断返回；再在 `main()` 增加一次调用并打印 `errors.Is`（打印需按输出注释规则）。

参考答案 2（可运行做法）：  
运行：
```bash
go run ./cmd/day01_07_slices_maps
# Output: ... iter: a=1 b=2 c=3 ... (输出可能变化/不固定：map 迭代顺序由运行时随机化)
```

---

## References
- 官方：Go 文档入口 https://go.dev/doc/
- 官方：Go Spec（含分号插入）https://go.dev/ref/spec
- 官方：builtin（make/len/cap）https://pkg.go.dev/builtin
- 官方：encoding/json（nil slice vs empty slice）https://pkg.go.dev/encoding/json
