# go-study（Go 全栈学习仓库）

这个仓库用于“Go 全栈作品集路线”的学习与沉淀：**可运行代码 + 可复习笔记 + 外部记忆（防 token 膨胀）+ 可重开窗口继续**。

## 快速定位（重开窗口先看这里）

- 学习进度索引（外部记忆 / Context Index）  
  - `/Users/zhang/Desktop/go-study/codex/notes/progress.md`
- 当天笔记目录  
  - `/Users/zhang/Desktop/go-study/codex/notes`
- Go module（所有 `go run/go test/go mod` 都在这里执行）  
  - `/Users/zhang/Desktop/go-study/codex/go-learning`
  - 代码入口索引：`/Users/zhang/Desktop/go-study/codex/go-learning/README.md`
- 教学 skill（固化流程，支持“查看学习进度/继续学习/开始学习”触发）  
  - `/Users/zhang/Desktop/go-study/codex/go-fullstack-coach/SKILL.md`

## 运行方式（最常用）

```bash
# 进入 Go module
cd /Users/zhang/Desktop/go-study/codex/go-learning

# 运行某个示例（当前保留的知识点示例都在 cmd/kp/*）
go run ./cmd/kp/basic-types
```

## 外部记忆（为什么要有 progress/glossary/patterns/pitfalls）

长对话很容易 token 膨胀。这里把“必须长期保留的上下文”写进文件：

- `notes/progress.md`：你学到哪了、下一步是什么（新对话只需读这份）
- `notes/glossary.md`：术语表（避免每次重新解释）
- `notes/patterns.md`：可复用工程套路（HTTP/DB/错误码等）
- `notes/pitfalls.md`：踩坑清单（避免反复踩同一类坑）

## README 同步规则（强制）

只要发生以下任意变化，就需要同步更新本 README（以及相关子 README）：

- 顶层目录结构变化（例如新增/移动 `go-learning/`、`notes/`、`go-fullstack-coach/`）
- 学习流程/触发口令变化（例如 skill 增加新指令、产物路径变更）
- 外部记忆文件清单变化（progress/glossary/patterns/pitfalls 增减）
- 关键运行方式变化（例如 Go module 迁移、cmd 归档规则调整）
