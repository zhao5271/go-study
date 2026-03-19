# Go 全栈学习进度（外部记忆 / Context Index）

> 目的：解决对话 token 膨胀。新对话只需要读这份索引 + 当天笔记即可。

## 当前状态
- 目标：从 Vue3/TypeScript/Node 转 Go 全栈（通用 API + MySQL）
- 教学规则（当前）：每次回复优先只讲 1–3 个知识点；讲解步骤不死板（不强制 A–I 模板），但要把“定义/为什么重要/关键坑与取舍/最小可运行代码/怎么运行+典型输出”讲清楚；所有输出点必须标注典型输出/不确定性原因；NotebookLM 优先，不足再上网补齐并写 References；不自动转 HTML；`go test` 非强制（你要求才跑）
- 当前进度：已完成 Day 05.2（MySQL + 列表 API + 错误码草案），准备进入 Day 06（Gin 工程化落地）
- 点播笔记：`notes/go/kp/01-变量与常量（作用域、_、iota）.md`（变量/常量：作用域、`_`、`iota`）
- 点播笔记：`notes/go/kp/02-基础类型（转换、格式化、表达式）.md`（类型转换、fmt/strconv、运算符与表达式）

## 已完成
| Day | 主题 | 产物 | 复习要点 |
|---|---|---|---|
| Day 01 | 语法地基：package/import、变量/零值、函数/error、控制流、slice/map | 笔记：`notes/go/day01-go-syntax-basics.md`；代码：`go-learning/cmd/day01/01_packages_import` ~ `go-learning/cmd/day01/07_slices_maps` | 零值、无隐式转换、`(value, error)`、slice 共享底层数组、nil/empty、map 顺序不稳定 |
| Day 02 | 函数与错误处理：wrap、errors.Is/As、shadowing、defer、panic/recover、table-driven tests | 笔记：`notes/go/day02-functions-errors.md`；代码：`go-learning/cmd/day02/01_returns` ~ `go-learning/cmd/day02/05_panic_recover`；测试：`go-learning/internal/day02/users/users_test.go` | `%w` vs `%v`、sentinel vs typed error、shadowing、defer 三条规则、panic 边界、table-driven tests |
| Day 03.1 | struct + 方法（值/指针接收者）+ API optional 字段（指针字段） | 笔记：`notes/go/day03-struct-methods-basics.md`；代码：`go-learning/cmd/day03/01a_struct_zero`、`go-learning/cmd/day03/01b_methods_receivers`、`go-learning/cmd/day03/01c_json_optional_fields` | zero value、值/指针接收者取舍、PATCH 区分“缺失 vs 显式零值”（`*T` + `omitempty`） |
| Day 03.2 | 接口（interface）+ 隐式实现 + embedding：Repo/Service 分层与装饰器 | 笔记：`notes/go/day03-2-interfaces-embedding.md`；代码：`go-learning/cmd/day03/02_interfaces_embedding`（含 ex1/ex2） | interface 贴业务语义、service 依赖接口、repo 可替换、embedding 装饰器（log/metrics） |
| Day 04.1 | net/http 打底：ServeMux/Handler、统一 JSON 响应与错误码、分页 query | 笔记：`notes/go/day04-net-http-basics.md`；代码：`go-learning/cmd/day04/01a_http_mux`、`go-learning/cmd/day04/01b_json_errors` | /api/v1 版本化、405/400 一致 JSON、HTTP status vs 业务 code、分页边界（size 限制） |
| Day 04.2 | 抽可复用 HTTP 工具包：`writeJSON/writeError/parsePageSize` | 笔记：`notes/go/day04-2-httpkit.md`；代码：`go-learning/internal/httpkit/*`（被 day04/day05 示例复用） | 统一响应出口、统一分页校验、减少复制粘贴导致的不一致 |
| Day 05.1 | MySQL 入门优先：Docker Compose + init.sql + Go 连接 + 事务回滚 | 笔记：`notes/go/day05-1-mysql-docker-crud.md`；compose：`go-learning/infra/mysql/docker-compose.yml`；代码：`go-learning/cmd/day05/01_mysql_crud`、`go-learning/cmd/day05/02_mysql_tx_rollback` | 3307 端口映射、init.sql 只首次执行、DSN parseTime、事务原子性（回滚验证 0/0） |
| Day 05.2 | 把 MySQL 接回 API：ListUsers 分页 + 错误码表草案 + net/http handler | 笔记：`notes/go/day05-2-list-users-errors-http.md`；代码：`go-learning/cmd/day05/03_list_users_sql`、`go-learning/cmd/day05/04_http_users_list` | LIMIT/OFFSET + COUNT、search LIKE 的索引取舍、DB→业务码→HTTP 映射、404/405 一致 JSON |

## 下一步（建议）
- Day 06：Gin 工程化落地（路由分组/中间件/结构化日志/错误码抽包），并把 users list 接到 Gin
