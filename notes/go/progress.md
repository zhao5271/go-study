# Go 全栈学习进度（外部记忆 / Context Index）

> 目的：解决对话 token 膨胀。新对话只需要读这份索引 + 当天笔记即可。

## 当前状态
- 目标：从 Vue3/TypeScript/Node 转 Go 全栈（通用 API + MySQL）
- 教学规则：每次输出严格 1)~5) 结构；示例代码 `fmt.Print*` 必须同行写典型输出注释；NotebookLM 优先，不足再上网补齐并写 References；不自动转 HTML；`go test` 暂不强制
- 进行中：Day 03（struct / 方法 / 接口）— 已完成 Day 03.1（见 `notes/go/day03-struct-methods-basics.md`）

## 已完成
| Day | 主题 | 产物 | 复习要点 |
|---|---|---|---|
| Day 01 | 语法地基：package/import、变量/零值、函数/error、控制流、slice/map | 笔记：`notes/go/day01-go-syntax-basics.md`；代码：`go-learning/cmd/day01_01_packages_import` ~ `go-learning/cmd/day01_07_slices_maps` | 零值、无隐式转换、`(value, error)`、slice 共享底层数组、nil/empty、map 顺序不稳定 |
| Day 02 | 函数与错误处理：wrap、errors.Is/As、shadowing、defer、panic/recover、table-driven tests | 笔记：`notes/go/day02-functions-errors.md`；代码：`go-learning/cmd/day02_01_returns` ~ `go-learning/cmd/day02_05_panic_recover`；测试：`go-learning/internal/day02/users/users_test.go` | `%w` vs `%v`、sentinel vs typed error、shadowing、defer 三条规则、panic 边界、table-driven tests |
| Day 04 | net/http 打底：ServeMux/Handler、统一 JSON 响应与错误码、分页 query | 笔记：`notes/go/day04-net-http-basics.md`；代码：`go-learning/cmd/day04_01a_http_mux`、`go-learning/cmd/day04_01b_json_errors` | /api/v1 版本化、405/400 一致 JSON、HTTP status vs 业务 code、分页边界（size 限制） |
| Day 05.1 | MySQL 入门优先：Docker Compose + init.sql + Go 连接 + 事务回滚 | 笔记：`notes/go/day05-1-mysql-docker-crud.md`；compose：`go-learning/infra/mysql/docker-compose.yml`；代码：`go-learning/cmd/day05_01_mysql_crud`、`go-learning/cmd/day05_02_mysql_tx_rollback` | 3307 端口映射、init.sql 只首次执行、DSN parseTime、事务原子性（回滚验证 0/0） |
| Day 05.2 | 把 MySQL 接回 API：ListUsers 分页 + 错误码表草案 + net/http handler | 笔记：`notes/go/day05-2-list-users-errors-http.md`；代码：`go-learning/cmd/day05_03_list_users_sql`、`go-learning/cmd/day05_04_http_users_list` | LIMIT/OFFSET + COUNT、search LIKE 的索引取舍、DB→业务码→HTTP 映射、404/405 一致 JSON |

## 下一步（建议）
- Day 06：Gin 工程化落地（路由分组/中间件/结构化日志/错误码抽包），并把 users list 接到 Gin
