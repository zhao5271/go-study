# Go 全栈学习进度（外部记忆 / Context Index）

> 目的：解决对话 token 膨胀。新对话只需要读这份索引 + 当天笔记即可。

## 当前状态
- 目标：从 Vue3/TypeScript/Node 转 Go 全栈（通用 API + MySQL）
- 教学规则：每次输出严格 1)~5) 结构；示例代码 `fmt.Print*` 必须同行写典型输出注释；NotebookLM 优先，不足再上网补齐并写 References；不自动转 HTML；`go test` 暂不强制

## 已完成
| Day | 主题 | 产物 | 复习要点 |
|---|---|---|---|
| Day 01 | 语法地基：package/import、变量/零值、函数/error、控制流、slice/map | 笔记：`notes/go/day01-go-syntax-basics.md`；代码：`go-learning/cmd/day01_01_packages_import` ~ `go-learning/cmd/day01_07_slices_maps` | 零值、无隐式转换、`(value, error)`、slice 共享底层数组、nil/empty、map 顺序不稳定 |
| Day 02 | 函数与错误处理：wrap、errors.Is/As、shadowing、defer、panic/recover、table-driven tests | 笔记：`notes/go/day02-functions-errors.md`；代码：`go-learning/cmd/day02_01_returns` ~ `go-learning/cmd/day02_05_panic_recover`；测试：`go-learning/internal/day02/users/users_test.go` | `%w` vs `%v`、sentinel vs typed error、shadowing、defer 三条规则、panic 边界、table-driven tests |

## 下一步（建议）
- Day 03：`struct` / 方法 / 接口（对照 TS interface）→ 为后续 HTTP 分层与 MySQL repository 奠基

