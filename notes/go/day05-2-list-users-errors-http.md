# Day 05.2：把 MySQL 接回后台管理 API（用户列表分页）+ 错误模型（DB → 业务码 → HTTP）

> 本课只讲 2 个知识点：  
> 1) `ListUsers(page,size,search)`：SQL + 索引思维  
> 2) 错误模型：DB error → 业务错误码 → HTTP 语义（统一 JSON）

---

## 知识点 1：用户列表分页查询（`LIMIT/OFFSET` + 可选搜索）

### B. 一句话定义
把“列表分页 + 可选搜索”封装成一个函数 `ListUsers(page,size,search)`，输出 items + total，并对 page/size 做严格边界校验。

### C. 为什么重要（不做会怎样）
后台管理最常见接口就是列表；没有分页会导致一次拉全表，性能/体验/成本都会爆炸；没有边界校验会被恶意/误用参数拖垮服务。

### D. 重难点拆解（2–4 条）
1) **分页边界**：`page>=1`，`1<=size<=100`（示例沿用此前 API 约定）。  
2) **total 的计算**：通常要 `COUNT(*)` 单独查一次（最小方案）；后面再讲优化。  
3) **搜索与索引**：`LIKE '%x%'` 会让索引很难用（最小示例先跑通，后续再升级搜索方案）。

### E. 业务场景落地（后台管理 API）
- 用户管理页：按关键字（email/name）搜索 + 分页展示。

### F. 代码示例（最小可运行）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day05/03_list_users_sql/main.go`

### G. 怎么运行（命令 + 预期现象）
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
MYSQL_DSN="app:app@tcp(127.0.0.1:3307)/go_admin?parseTime=true" go run ./cmd/day05/03_list_users_sql
# Output: total=2 (输出可能变化/不固定：取决于数据与 search)
# Output: user: id=1 email=admin@example.com name=Admin role=admin
```

### H. 练习题（1–3 题）
练习 1：把 page/size/search 作为命令行参数传入（比如 `-page -size -search`）  
- 验收标准：
  - `-page 0` 会报错（你可以打印错误并退出）
  - `-size 101` 会报错
  - `-search ""` 能返回全量（但仍分页）

### I. 参考答案
参考答案 1：你可以用标准库 `flag` 包实现（下一课我们会把“参数解析 + 校验”抽成可复用层）。

---

## 知识点 2：错误模型（DB error → 业务错误码 → HTTP status）+ 统一 JSON 响应

### B. 一句话定义
把不同失败原因稳定地映射成：HTTP status（协议语义）+ 业务 code（前端/告警可聚合）+ 统一 JSON 响应结构。

### C. 为什么重要（不做会怎样）
- 前端很难稳定处理（提示/重试/埋点/灰度）；  
- 运维/你自己难以用日志快速定位问题（同类错误无法按 code 聚合）。

### D. 重难点拆解（2–4 条）
1) **400 vs 500**：参数问题是 400（调用方负责修）；DB 挂了是 500（服务端负责）。  
2) **错误码表要“稳定”**：一旦前端/监控依赖了 code，就不能随便改。  
3) **404/405 也要一致 JSON**：否则调用方还是要分支处理。

### E. 业务场景落地（后台管理 API）
- `GET /api/v1/users?page=1&size=20&search=a`：用户列表页的核心接口（后面再加鉴权/角色权限）。

### F. 代码示例（最小可运行）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day05/04_http_users_list/main.go`

错误码表（草案）：
| HTTP | code | message | 场景 |
|---:|---:|---|---|
| 200 | 0 | OK | 成功 |
| 400 | 10002 | INVALID_QUERY | page/size 非法 |
| 404 | 10004 | NOT_FOUND | 路由不存在 |
| 405 | 10001 | METHOD_NOT_ALLOWED | method 不支持 |
| 500 | 20001 | DB_ERROR | DB 连接/查询错误 |
| 500 | 20002 | INTERNAL_ERROR | 其他未分类错误 |

### G. 怎么运行（命令 + 预期现象）
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
PORT=18080 MYSQL_DSN="app:app@tcp(127.0.0.1:3307)/go_admin?parseTime=true" go run ./cmd/day05/04_http_users_list
# Output: 2006/01/02 15:04:05 listening on :18080 (输出可能变化/不固定：包含时间戳)
```

验证（另开终端）：
```bash
curl -s "http://localhost:18080/api/v1/users?page=1&size=2&search=a"
# Output: {"code":0,"message":"OK","data":{"items":[...],"page":1,"size":2,"total":2}} (输出可能变化/不固定：items/total 随数据变化)

curl -s "http://localhost:18080/api/v1/users?page=0"
# Output: {"code":10002,"message":"INVALID_QUERY"}

curl -s -X POST "http://localhost:18080/api/v1/users"
# Output: {"code":10001,"message":"METHOD_NOT_ALLOWED"}

curl -s "http://localhost:18080/api/v1/not-exist"
# Output: {"code":10004,"message":"NOT_FOUND"}
```

### H. 练习题（1–3 题）
练习 1：为 `/api/v1/users` 增加 `search` 的长度限制（例如 <= 50）  
- 验收标准：超长 search 返回 400 + `INVALID_QUERY`

练习 2：把 `DB_ERROR` 的场景做一个“可复现演示”  
- 验收标准：停掉 MySQL 容器后再 curl，返回 500 + `DB_ERROR`

### I. 参考答案
参考答案 1：在 handler 里 `strings.TrimSpace(search)` 后判断 `len(search) > 50`，直接 `httpkit.WriteError(400,10002,"INVALID_QUERY")`（或你的同名封装）。

参考答案 2：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning/infra/mysql
docker compose stop
# Output: ... (输出可能变化/不固定：取决于 Docker 版本)

curl -s "http://localhost:18080/api/v1/users?page=1&size=2"
# Output: {"code":20001,"message":"DB_ERROR"}
```

---

## References
- 官方：Go `net/http` https://pkg.go.dev/net/http （官方）
- 官方：Go `database/sql` https://pkg.go.dev/database/sql （官方）
- 驱动：go-sql-driver/mysql https://pkg.go.dev/github.com/go-sql-driver/mysql （官方/准官方）
- 规范参考：`api-design-principles`（资源命名、方法语义、分页参数、错误响应一致性）
