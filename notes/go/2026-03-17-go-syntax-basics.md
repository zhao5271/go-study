# Go 语法基础（TS/Vue/Node 转 Go）- Lesson 01

日期：2026-03-17

## 今日目标
- 用一份**可运行的最小程序**把 Go 的“语法地基”搭起来：`package/import`、变量与零值、函数多返回值与 `error`、`if` 初始化语句、`for/switch`、`slice/map`、`nil vs empty`。

## 背景对照（TS/Node → Go）
- Go 的设计倾向：**简单、可读、编译期约束强**（未使用的变量/导入直接编译失败），换取大型工程的可维护性与编译速度。
- Go 没有 `undefined`：声明就有**零值**，减少“没初始化就用”的问题。
- Go 常用 `(value, error)` 多返回值替代 try/catch，把错误处理显式化。

## 关键结论
- `var`：适合包级、需要零值、需要显式类型；`:=`：只在函数内、快速声明 + 类型推导。
- `if err := ...; err != nil {}`：把“临时变量作用域”锁在 if 内，避免污染外部。
- `for` 是唯一循环关键字；`switch` 默认 `break`，只有 `fallthrough` 才会继续下一个 case。
- `slice` 是“底层数组的视图”，可能共享底层数组；想隔离就 `copy`。
- `nil slice` 可以 `append`；`nil map` 不能写入（会 panic），必须 `make(map[...])`。

## 代码清单（可运行）
- 入口演示：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/lesson01/main.go:1`
- 错误处理示例函数：`/Users/zhang/Desktop/go-study/codex/go-learning/internal/basics/divide.go:1`
- 测试：`/Users/zhang/Desktop/go-study/codex/go-learning/internal/basics/divide_test.go:1`

运行：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go test ./...
go run ./cmd/lesson01
```

## 常见坑（结合 TS/Node 习惯）
- `:=` 变量遮蔽（shadowing）：在 `if`/`for` 里不小心创建了同名新变量，外层变量没被更新。
- `range` 拿到的是**值的拷贝**（对 slice 元素），想修改原数据要用索引 `s[i] = ...`。
- `map` 遍历顺序不稳定：不要依赖输出顺序（调试打印要心里有数）。
- `slice` 共享底层数组：对“子切片”的修改可能影响原切片；append 触发扩容后共享关系又会变化。

## 面试问法（你可以用来复述）
- Go 为什么把错误处理设计成 `error` 返回值而不是异常？
- `slice` 的 `len`/`cap` 各表示什么？什么时候会共享底层数组？
- `nil slice` vs `empty slice` 的区别是什么？什么时候会影响 JSON 输出？

## 下一步
- Lesson 02：类型系统与“组合类型”——`struct`、方法、接口（用 TS 的 interface/duck typing 做对照），并为后续 HTTP API 做铺垫。

