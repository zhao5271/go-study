# 知识点点播：变量与常量（作用域 / 匿名变量 `_` / `iota`）

> 目标：把“变量/常量的可见范围 + 忽略值 + iota 常量枚举”一次讲透，并能直接跑一个最小示例。  
> 贯穿项目落地：后台管理 API 里会经常遇到“临时变量/错误变量的作用域”“忽略不需要的返回值”“用常量表达稳定的业务枚举/错误码/角色”等。

## 关联复习（用到你已学过的点）
- Day01：变量/零值/短变量声明 `:=`：`/Users/zhang/Desktop/go-study/codex/notes/go/day01-go-syntax-basics.md:1`
- Day02：shadowing（`:=` 可能创建新变量遮蔽外层）：`/Users/zhang/Desktop/go-study/codex/notes/go/day02-functions-errors.md:1`

## 说明：资料来源
- NotebookLM 本次查询失败（`net::ERR_CONNECTION_CLOSED`），本节结论以 Go 官方文档为准（见文末 References）。

---

## 1) 作用域（scope）+ block scope + `:=` 的“遮蔽（shadowing）”

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

可运行示例（本节所有示例共用一个入口）：
- 代码：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/kp/vars-const-scope-iota/main.go:1`
- 运行：见下方“怎么运行”

---

## 2) 匿名变量 `_`（blank identifier）：显式“我不需要这个值”

一句话：`_` 是一个“**丢弃桶**”，把你不需要的返回值/循环变量明确丢掉。

为什么重要（后台管理 API 场景）：
- 标准库/你自己的函数经常多返回值：`(value, error)`、`(rows, err)`、`(affected, err)`。  
  你只想要其中一个，就用 `_` 表达“我刻意忽略”，并通过编译器的“未使用即报错”保持代码干净。

关键坑：
1) `_` 不是“随便写写”，它会让你**永久失去那个值**；别把重要的错误/返回值用 `_` 吞了。  
2) `range` 里如果你不用 index/value，优先用 `_`，避免写一个没意义的变量名。

可运行示例：同上（`demoBlankIdentifier()`）。

---

## 3) `const` + `iota`：写“稳定枚举/位标志”的最省心方式

一句话：`iota` 只在 `const (...)` 块里出现，**从 0 开始、每行 +1**，常用于枚举/位标志。

为什么重要（后台管理 API 场景）：
- 角色（Admin/Viewer）、状态（Enabled/Disabled）、错误码类别、审计动作类型……这些都是**稳定、可复用、可搜索**的常量，适合用 `const` 表达，而不是散落的 magic number/string。

关键坑：
1) `iota` **按“行”递增**（不是按“写了 iota 才递增”）：哪怕某行你写了显式值，它也会继续 +1。  
2) `iota` **每个 const 块都会重置为 0**：跨文件/跨块不共享计数。

可运行示例：同上（`demoIota()`）。

---

## 怎么运行 + 预期现象

```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/kp/vars-const-scope-iota
```

你会看到三段输出（顺序固定）：
- scope/shadowing：inner/outer 的值不同
- blank identifier：忽略一个返回值仍能正常运行
- iota：枚举值与位标志数值打印出来

---

## 自检/练习（不额外加目录，直接改本文件即可）

练习 1：把 `demoScope()` 里 `if true { x := 2 ... }` 改成 `x = 2`  
- 验收：外层 `x` 也变成 2（输出从 `outer x=1` 变成 `outer x=2`）

练习 2：在 `demoBlankIdentifier()` 里把 `_` 改成变量名 `rem` 并打印它  
- 验收：你能看到余数（例如 10%3=1）

练习 3：给 `Role` 增加一个角色 `RoleGuest`，并确认它的值是递增的  
- 验收：打印出来是连续的（1,2,3…）

---

## References
- 官方：Go Spec（Constants / Iota / Scopes / Blank identifier）https://go.dev/ref/spec  
- 官方：Effective Go（Constants / iota 常见用法）https://go.dev/doc/effective_go  
- 官方：A Tour of Go（Constants）https://go.dev/tour/basics/15

