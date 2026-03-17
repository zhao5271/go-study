# Go 全栈学习笔记 - Lesson 01.1：变量声明、零值、类型转换

> 目标：你是 Vue3/TypeScript/Node 工程师，第一步先把 Go 的“变量与零值体系”建立起来，这是后面写 API/接 MySQL/做并发的地基。

## 1) 知识讲解：概念 → 为什么（设计动机/取舍）

### 1.1 `var` 声明 + 零值（Zero Value）
**概念**
- `var x T` 声明一个变量 `x`，类型是 `T`。
- Go **没有 `undefined`**：变量声明后立刻拥有**零值**：
  - 数字：`0`
  - 字符串：`""`
  - 布尔：`false`
  - 指针/切片/map/函数/chan/interface：`nil`

**为什么**
- 取舍：牺牲了“未定义状态”的表达力，换来工程稳定性（减少“没初始化就用”的问题）。

**对照 TS/Node**
- TS/JS 常见坑：`undefined`/`null` 混用、对象字段不存在导致运行期异常。
- Go 更偏“编译期约束 + 运行期稳定”，把很多问题提前到编译期/零值语义上解决。

#### 可运行例子（紧跟）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/lesson01_01a_var_zero/main.go:1`

运行：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/lesson01_01a_var_zero
```

代码（可直接阅读）：
```go
package main

import "fmt"

func main() {
	fmt.Println("== 01.1A var + zero values ==") // Output: == 01.1A var + zero values ==

	var count int
	var name string
	var ok bool
	fmt.Printf("count=%d name=%q ok=%v\n", count, name, ok) // Output: count=0 name="" ok=false
}
```

你会看到代码里每个 `fmt.Print*` 的同一行注释都写了典型输出。

---

### 1.2 `:=` 短变量声明（只在函数内）
**概念**
- `x := 42` = 声明 + 赋值 + 类型推导（类型一旦确定就固定）。

**为什么**
- 取舍：让局部代码更短更快写；但仍然是**静态强类型**，不是 JS 那种“随时换类型”。

**对照 TS/Node**
- 很像 TS 的 `let x = 42` 的推导体验，但 Go 的类型不可变更，且不做隐式类型转换。

#### 可运行例子（紧跟）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/lesson01_01b_short_decl/main.go:1`

运行：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/lesson01_01b_short_decl
```

代码（可直接阅读）：
```go
package main

import "fmt"

func main() {
	fmt.Println("== 01.1B := short declaration ==") // Output: == 01.1B := short declaration ==

	x := 42
	fmt.Printf("x=%d type=%T\n", x, x) // Output: x=42 type=int
}
```

---

### 1.3 数字类型不隐式转换（`int` / `int64` 等）
**概念**
- `int` 不能直接加 `int64`，必须显式转换：`int64(x) + big`

**为什么**
- 取舍：多写一点类型转换，换来“不会悄悄溢出/丢精度”的可控性。

**对照 TS/Node**
- JS 里数字基本都是 `number`（浮点），很多边界问题在运行期才暴露。
- Go 把这类问题通过类型系统更早暴露出来。

#### 可运行例子（紧跟）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/lesson01_01c_conversion/main.go:1`

运行：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/lesson01_01c_conversion
```

代码（可直接阅读）：
```go
package main

import "fmt"

func main() {
	fmt.Println("== 01.1C explicit conversion (int vs int64) ==") // Output: == 01.1C explicit conversion (int vs int64) ==

	x := 42
	var big int64 = 1
	sum := int64(x) + big
	fmt.Printf("int64(x)+big=%d type=%T\n", sum, sum) // Output: int64(x)+big=43 type=int64
}
```

---

### 1.4 常量 `const`
**概念**
- `const` 声明常量，很多常量在需要时才“落地成具体类型”（未类型化常量）。

**为什么**
- 取舍：让常量在不同上下文里更好用（例如参与不同类型的表达式时更灵活），同时保持编译期可计算。

**对照 TS/Node**
- TS 的 `const` 更像“不可重新赋值的变量”；Go 的 `const` 更强调“编译期常量”属性。

#### 可运行例子（紧跟）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/lesson01_01d_constants/main.go:1`

运行：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/lesson01_01d_constants
```

代码（可直接阅读）：
```go
package main

import "fmt"

func main() {
	fmt.Println("== 01.1D constants ==") // Output: == 01.1D constants ==

	const pi = 3.14159
	fmt.Printf("pi=%.2f type=%T\n", pi, pi) // Output: pi=3.14 type=float64
}
```

## 2) 常见坑（结合 TS/Node 习惯对照）
- **`:=` 只能在函数内**：包级变量只能用 `var`。
- **未使用就报错**：声明变量但没用会编译失败（Go 的工程“洁癖”）。
- **`:=` 变量遮蔽（shadowing）**：在 `if/for` 里写 `:=` 可能创建同名新变量，外层变量没更新（后续我们会专门讲）。
- **别指望隐式转换**：`int` / `int64` / `float64` 等都要显式转换。

## 3) 工程用法/最佳实践（真实 API 项目怎么用）
- 局部变量默认优先 `:=`（短、清晰），需要指定类型/包级变量才用 `var`。
- “类型宽度重要”的场景（时间戳、字节大小、数据库字段）更建议显式类型。
- 养成习惯：代码写完第一件事跑 `gofmt`（后面配合 editor 自动格式化）。

## 4) 练习策略（练习=运用示例，给完整参考实现）
练习就用这份参考实现：
- 参考实现：
  - `go run ./cmd/lesson01_01a_var_zero`（零值）
  - `go run ./cmd/lesson01_01b_short_decl`（短声明）
  - `go run ./cmd/lesson01_01c_conversion`（显式转换）
  - `go run ./cmd/lesson01_01d_constants`（常量）

你可以做 2 个“无痛改造”（不额外加新文件）：
1) 新增一个变量 `var ratio float64`，打印它的零值与类型（照着现有输出注释规则写）。
2) 故意写一行 `sum2 := x + big` 观察编译报错，再改成正确的显式转换。

下一步建议：Lesson 01.2 进入“函数 + `(value, error)` + `if err != nil`”，这是 Go 后端最核心的工程习惯之一。
