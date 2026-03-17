# Day 05.1：MySQL 入门优先（Docker 可复现）+ Go 连接（database/sql）+ 事务回滚

默认贯穿项目：后台管理 API（RBAC + 登录鉴权 + 列表分页检索 + CRUD + 审计日志 + Docker 部署）

> 你当前画像：目标=作品集交付；Docker=能；MySQL=没有基础。  
> 本课目标：先把“本地可复现 MySQL + 最小 CRUD + 事务回滚”跑通，再接回 HTTP API。

---

## 知识点 1：用 Docker Compose 一键启动 MySQL（带初始化表与种子数据）

### B. 一句话定义
用 `docker compose` 把 MySQL 当作“项目依赖”启动起来，并通过 init SQL 自动建表和插入种子数据，做到新机器也能一条命令复现。

### C. 为什么重要（不做会怎样）
- 作品集项目必须能复现：面试官/同事拉下来要能跑；否则“我本地能跑”没有说服力。
- 没有初始化脚本会导致环境漂移：你写的代码依赖某些表/数据，但别人机器没有。

### D. 重难点拆解（2–4 条）
1) **端口占用**：本仓库把 MySQL 映射到 `3307`，避免与你本机可能存在的 `3306` 冲突。  
2) **初始化脚本只在首次创建数据卷时执行**：如果你改了 `init.sql`，需要清理 volume 才会重新跑。  
3) **字符集**：默认 `utf8mb4`，避免中文/emoji 乱码。

### E. 业务场景落地（后台管理 API）
- 先落一张 `users` 表和一张 `audit_logs` 表：后面做“创建用户 + 写审计日志”必须保证一致性（事务）。

### F. 代码示例
文件：
- `/Users/zhang/Desktop/go-study/codex/go-learning/infra/mysql/docker-compose.yml`
- `/Users/zhang/Desktop/go-study/codex/go-learning/infra/mysql/init.sql`

### G. 怎么运行（命令 + 预期现象）
启动 MySQL：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning/infra/mysql
docker compose up -d
# Output: ... Started ... (输出可能变化/不固定：取决于 Docker 版本与拉取进度)
```

检查容器健康（可选）：
```bash
docker ps --format 'table {{.Names}}\t{{.Status}}'
# Output: go-learning-mysql  Up ... (healthy) (输出可能变化/不固定：状态/时间会变化)
```

进入容器并查询种子数据：
```bash
docker exec -it go-learning-mysql mysql -uapp -papp -D go_admin -e "SELECT id,email,name,role FROM users ORDER BY id;"
# Output: id  email               name   role
# Output: 1   admin@example.com   Admin  admin
# Output: 2   alice@example.com   Alice  editor
# Output: 3   bob@example.com     Bob    viewer
```

停止并清理（谨慎：会删数据卷，重新执行 init.sql）：
```bash
docker compose down -v
# Output: ... removed ... (输出可能变化/不固定：取决于 Docker 版本)
```

### H. 练习题（1–3 题）
练习 1：把 `users.role` 的默认值改成 `'viewer'`（已经是）并验证新插入用户不传 role 也会是 viewer  
- 验收标准：执行一条 INSERT（不传 role），SELECT 出 role=viewer

### I. 参考答案
参考答案 1（可运行）：
```bash
docker exec -it go-learning-mysql mysql -uapp -papp -D go_admin -e "INSERT INTO users (email,name) VALUES ('role_default@example.com','Role Default'); SELECT email,role FROM users WHERE email='role_default@example.com';"
# Output: role_default@example.com  viewer
```

---

## 知识点 2：Go 连接 MySQL（database/sql）并做最小 CRUD（插入 + 查询）

### B. 一句话定义
用标准库 `database/sql` + MySQL 驱动连接数据库，通过 `Exec/Query` 做最小可运行的插入与查询。

### C. 为什么重要（不做会怎样）
你后面做“登录/列表/权限”都离不开数据库；如果连“连接 + 查询 + 插入”都不稳定，工程交付会一直卡住。

### D. 重难点拆解（2–4 条）
1) **DSN 参数**：建议 `parseTime=true`，否则 `DATETIME` 扫描到 `time.Time` 会出问题。  
2) **超时**：DB 操作必须可控（先用 `context.WithTimeout`）；否则慢查询会拖垮服务。  
3) **驱动选择**：示例用 `github.com/go-sql-driver/mysql`（事实标准）；只为 MySQL 连接服务，不引入 ORM。

### E. 业务场景落地（后台管理 API）
- “创建用户”：插入 `users`；后面会加密码哈希、唯一约束冲突处理（409）。

### F. 代码示例（最小可运行）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day05_01_mysql_crud/main.go`

### G. 怎么运行（命令 + 预期现象）
先确保 MySQL 已启动（知识点 1）。

安装依赖（在 `go-learning` 目录执行）：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go get github.com/go-sql-driver/mysql
# Output: go: added github.com/go-sql-driver/mysql ... (输出可能变化/不固定：版本号会变)
```

运行：
```bash
MYSQL_DSN="app:app@tcp(127.0.0.1:3307)/go_admin?parseTime=true" go run ./cmd/day05_01_mysql_crud
# Output: db ping ok
# Output: inserted user id=4 (输出可能变化/不固定：自增 id 取决于现有数据)
# Output: user: id=1 email=admin@example.com name=Admin role=admin
```

### H. 练习题（1–3 题）
练习 1：把插入的 role 改成 `editor`，并查询验证  
- 验收标准：输出中能看到你插入的那条记录 role=editor（或你能用 SQL 查出来）

### I. 参考答案
参考答案 1：直接把 `role := "viewer"` 改成 `role := "editor"` 重新运行即可（输出的 role 会变化）。

---

## 知识点 3：事务（Transaction）= “多步写入要么全成功，要么全失败”

### B. 一句话定义
事务把多条 SQL 写操作包成一个原子单元：任何一步失败就回滚，避免“半成功”产生脏数据。

### C. 为什么重要（不做会怎样）
后台管理常见场景：**创建用户 + 写审计日志**。如果不事务，你可能“用户创建成功但审计日志失败”，或者反过来，数据一致性被破坏，排查困难。

### D. 重难点拆解（2–4 条）
1) **失败后必须回滚**：不回滚会占用连接/锁资源。  
2) **把真正的错误制造出来验证**：例如触发 UNIQUE 冲突，证明回滚真的生效。  
3) **事务里的查询/写入都要用 tx**：不要混用 db 和 tx（否则你以为在事务里，实际上不在）。

### E. 业务场景落地（后台管理 API）
- “创建用户（users）+ 记录操作（audit_logs）”必须同事务，保证一致性。

### F. 代码示例（最小可运行）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day05_02_mysql_tx_rollback/main.go`

### G. 怎么运行（命令 + 预期现象）
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
MYSQL_DSN="app:app@tcp(127.0.0.1:3307)/go_admin?parseTime=true" go run ./cmd/day05_02_mysql_tx_rollback
# Output: db ping ok
# Output: 2006/01/02 15:04:05 expected err=Error 1062 ... (输出可能变化/不固定：错误信息随驱动/版本变化)
# Output: rollback ok
# Output: users count=0
# Output: audit_logs count=0
```

### H. 练习题（1–3 题）
练习 1：把“制造错误”的步骤改成插入 `audit_logs` 时 `actor_user_id=0`（如果你后面给它加了 NOT NULL/外键约束）或插入超长 action（触发约束）  
- 验收标准：仍然能证明回滚后 `users count=0` 且 `audit_logs count=0`

### I. 参考答案
参考答案 1：本示例当前用“重复 email 触发 UNIQUE 冲突”来制造失败；你只要保证事务里出现任意错误并调用 rollback，即可得到同样的 0/0 输出。

---

## References
- 官方/社区：MySQL Docker 镜像初始化机制（`/docker-entrypoint-initdb.d`） https://hub.docker.com/_/mysql （官方）
- 官方：Go `database/sql` 包文档 https://pkg.go.dev/database/sql （官方）
- 驱动：go-sql-driver/mysql DSN 与 parseTime https://pkg.go.dev/github.com/go-sql-driver/mysql （官方/准官方）
- 规范参考：`postgresql-table-design`（NOT NULL/约束/索引/命名的思维框架，落地到 MySQL 语法）

