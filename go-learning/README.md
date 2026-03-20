# go-learning（Go 全栈学习代码仓库）

本目录是 Go module（见 `go.mod`）。所有 `go run / go test / go mod` 命令都应在这里执行：

```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
```

## 目录结构约定

- `cmd/`
- `cmd/kp/`：当前保留的知识点示例入口（按知识点 slug 归档；每个子目录都是一个 `package main`）

学习笔记在仓库根目录（上一级）：
- `/Users/zhang/Desktop/go-study/codex/notes`

## 如何运行示例（通用）

进入本目录后，直接运行某个示例目录：

```bash
go run ./cmd/kp/vars-const-scope-iota
go run ./cmd/kp/basic-types
go run ./cmd/kp/string-basics
```

## README 同步规则（强制）

只要发生以下任意变化，就需要同步更新本 README：

- `cmd/` 示例新增/移动/重命名
- 统一的运行方式、环境变量约定变更
