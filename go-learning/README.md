# go-learning（Go 全栈学习代码仓库）

本目录是 Go module（见 `go.mod`）。所有 `go run / go test / go mod` 命令都应在这里执行：

```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
```

## 目录结构约定

- `cmd/`
  - `cmd/day01/`、`cmd/day02/`…：每天的可运行示例入口（每个子目录都是一个 `package main`）
  - `cmd/kp/`：知识点点播示例入口（按知识点 slug 归档）
  - `cmd/legacy/`：历史/废弃示例（保留做对照）
- `internal/`：可复用逻辑（供多个示例共享）
  - `internal/httpkit/`：统一 JSON 响应与 query 解析（`WriteJSON/WriteError/ParsePageSize`）
  - `internal/day02/users/`：Day02 错误处理示例的复用逻辑
- `infra/`
  - `infra/mysql/`：MySQL Docker Compose + 初始化 SQL（可复现环境）

学习笔记在仓库根目录（上一级）：
- `/Users/zhang/Desktop/go-study/codex/notes/go`

## 常用环境变量

- `PORT`：HTTP 示例监听端口（建议 `18080` 避免冲突）
- `MYSQL_DSN`：MySQL DSN（示例默认）：
  - `app:app@tcp(127.0.0.1:3307)/go_admin?parseTime=true`

## 如何运行示例（通用）

进入本目录后，直接运行某个示例目录：

```bash
go run ./cmd/day01/01_packages_import
go run ./cmd/day02/02_errors
PORT=18080 go run ./cmd/day04/01b_json_errors
```

## MySQL（Docker 可复现）

启动（需要 Docker Desktop/daemon 已启动）：

```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning/infra/mysql
docker compose up -d
```

验证种子数据：

```bash
docker exec -it go-learning-mysql mysql -uapp -papp -D go_admin -e "SELECT id,email,name,role FROM users ORDER BY id;"
```

停止（保留数据卷）：

```bash
docker compose stop
```

## README 同步规则（强制）

只要发生以下任意变化，就需要同步更新本 README：

- `cmd/` 示例新增/移动/重命名（尤其是 day 归档结构）
- `internal/` 新增可复用包或对外 API 变更
- `infra/` 目录结构/端口/初始化脚本变更
- 统一的运行方式、环境变量约定变更
