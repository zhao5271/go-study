# Go 全栈学习笔记 - Day 01：语法地基（TS/Vue/Node → Go）

> 今日目标：把 Go 的“包/导入、变量与零值、函数与错误、控制流、slice/map”打牢，并且形成工程化习惯（可运行代码 + 可复习笔记）。

## 1) 知识讲解：概念 → 为什么（设计动机/取舍）

### 1.1 `package` / `import` 与导出规则（首字母大小写）
**概念**
- Go 的组织单位是 **package**：目录通常对应一个包（工程上强调清晰依赖）。
- 可见性规则：**首字母大写 = 导出（public）**；小写 = 包内私有（package-private）。

**为什么**
- 取舍：不用 `public/private/export` 关键字，靠命名约定控制可见性，语法更轻。
- 编译期更严格：未使用的 import 会直接报错，逼你保持依赖干净。

**对照 TS/Node**
- TS 通过 `export` 控制模块边界；Go 则通过 package + 名字大小写控制对外 API。

### 1.2 `var` + 零值（Zero Value）
**概念**
- Go 没有 `undefined`：声明即初始化为零值（0/""/false/nil）。

**为什么**
- 取舍：牺牲“未定义状态”的表达力，换来工程稳定性（减少没初始化就用）。

**对照 TS/Node**
- TS/JS 里 `undefined/null` 分支多、运行期才爆；Go 把很多问题提前到编译期 + 零值语义。

### 1.3 `:=` 短变量声明（函数内）
**概念**
- `x := 42` = 声明 + 赋值 + 类型推导；只能在函数体内使用。

**为什么**
- 取舍：更快写更短，但仍然是静态强类型（类型确定后不会变）。

### 1.4 显式类型转换 + `const`
**概念**
- Go 不做隐式数值转换（`int`/`int64` 必须显式转换）。
- `const` 更强调“编译期常量”。

**为什么**
- 取舍：多写一点转换，换来精度/溢出等边界更可控。

### 1.5 函数 + `(value, error)`（替代 try/catch）
**概念**
- Go 常用多返回值：`(value, error)`，调用点显式检查 `err`。

**为什么**
- 取舍：控制流显式化，避免异常穿透边界“漏处理”。

### 1.6 控制流：`if` / `for` / `switch`
**概念**
- `if` 支持 init：`if x := ...; x > 0 {}`（让临时变量作用域更小）
- `for` 是唯一循环关键字（经典/while-like/range）
- `switch` 默认 `break`；只有 `fallthrough` 才继续下一个 case

### 1.7 `slice` / `map` 与 “nil vs empty”
**概念**
- `slice` 是“底层数组的视图”，可能共享底层数组（高性能但要理解引用语义）。
- `map` 必须 `make` 后才能写入；`nil map` 写会 panic。
- `nil slice` 与 `empty slice` 在 JSON 编码上可能不同：`null` vs `[]`。
- `map` 遍历顺序不保证稳定。

---

## 2) 示例驱动：每个知识点后立刻给一段可运行代码

统一进入项目：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
```

### 2.1 package/import + 导出规则
运行：
```bash
go run ./cmd/day01_01_packages_import
```

代码（全文）见：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day01_01_packages_import/main.go:1`

### 2.2 var + zero values
运行：
```bash
go run ./cmd/day01_02_vars_zero
```
代码（全文）见：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day01_02_vars_zero/main.go:1`

### 2.3 := short declaration
运行：
```bash
go run ./cmd/day01_03_short_decl
```
代码（全文）见：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day01_03_short_decl/main.go:1`

### 2.4 conversions + const
运行：
```bash
go run ./cmd/day01_04_conversions_const
```
代码（全文）见：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day01_04_conversions_const/main.go:1`

### 2.5 functions + (value, error)
运行：
```bash
go run ./cmd/day01_05_functions_error
```
代码（全文）见：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day01_05_functions_error/main.go:1`

### 2.6 if / for / switch
运行：
```bash
go run ./cmd/day01_06_if_for_switch
```
代码（全文）见：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day01_06_if_for_switch/main.go:1`

### 2.7 slices/maps + nil/empty + map order
运行：
```bash
go run ./cmd/day01_07_slices_maps
```
代码（全文）见：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day01_07_slices_maps/main.go:1`

---

## 3) 常见坑（结合 TS/Node 习惯对照）
- `{` 不能另起一行（Go 会自动插分号，容易编译报错）。
- 未使用的 import/变量会直接报错（不是“严格模式”，是 Go 的工程约束）。
- `:=` 只能在函数内；并且在 `if/for` 里容易造成 **shadowing**（变量遮蔽）。
- `int`/`int64` 不会自动转换；别指望像 JS 一样“都能算”。
- `slice` 可能共享底层数组：子切片改元素会影响原切片。
- `nil map` 写入会 panic，必须 `make`。
- `map` 遍历顺序不稳定：输出/逻辑不要依赖顺序（如需顺序，取 key 排序再遍历）。

## 4) 工程用法/最佳实践（真实 API 项目怎么落地）
- 写完就 `gofmt`（后面建议 IDE 开保存自动格式化）。
- 代码组织从 “`cmd/` 入口 + `internal/` 业务包” 开始，后续接 HTTP/MySQL 会更顺。
- 错误处理先把模式练熟：`v, err := ...; if err != nil { return ... }`。
- 用 `nil` 表达“未初始化/缺省”，但对外 API（尤其 JSON）要明确 `null` vs `[]` 的契约。

## 5) 练习策略（练习可直接作为“运用示例”，提供完整参考实现）
完整参考实现就是本笔记对应的 7 个可运行示例（直接跑即可）：
- `go run ./cmd/day01_01_packages_import`
- `go run ./cmd/day01_02_vars_zero`
- `go run ./cmd/day01_03_short_decl`
- `go run ./cmd/day01_04_conversions_const`
- `go run ./cmd/day01_05_functions_error`
- `go run ./cmd/day01_06_if_for_switch`
- `go run ./cmd/day01_07_slices_maps`

你可以做 3 个“就地改造”（不新增文件也行）：
1) 在 Day01.5 增加一个新的错误分支（比如输入负数），并保持输出注释规则。
2) 在 Day01.7 里把 `copy` 去掉，观察 `s3` 修改是否影响 `s2`，写下你的结论。
3) 在 map 遍历前先把 keys 拿出来排序（`sort.Strings`），让输出顺序固定。

## References
- 官方：https://go.dev/doc/ （Go 官方文档入口）
- 官方：https://go.dev/ref/spec （Go 语言规范：语法与分号插入规则等）
- 官方：https://go.dev/blog/defer-panic-and-recover （defer/panic/recover 的工作机制）
- 官方：https://pkg.go.dev/encoding/json （JSON 编码行为：slice 的 null/[] 等）
- 官方：https://pkg.go.dev/builtin （make/len/cap 等内建函数语义）
