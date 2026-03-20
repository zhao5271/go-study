---
type: index
domain: go
role: progress
status: active
tags:
  - go
  - index
  - progress
created: 2026-03-20
updated: 2026-03-20
---

# Go 全栈学习进度（外部记忆 / Context Index）

> 目的：解决对话 token 膨胀。新对话只需要读这份索引 + 当天笔记即可。

## 当前状态
- 目标：从 Vue3/TypeScript/Node 转 Go 全栈（通用 API + MySQL）
- 记忆方案与稳定规则：见 `notes/memory-strategy.md` 与仓库 `AGENTS.md`
- 当前目录已拍平：索引/记忆文件直接放在 `notes/*.md`，知识点笔记统一放在 `notes/kp/*.md`
- 历史学习进度仍以“已学到 Day 05.2（MySQL + 列表 API + 错误码草案）”为认知基线，但旧的 Day 系列笔记与示例已清理，不再作为当前仓库入口
- 最近同步：`notes/kp/01-03` 三篇基础知识点笔记已于 2026-03-20 按新 `knowledge-point-notes` skill 重整，并统一补齐官方 References / TL;DR / 练习与参考答案

## 当前保留资料
| 类型 | 主题 | 入口 | 配套代码 |
|---|---|---|---|
| KP 01 | 变量与常量（作用域、`_`、`iota`） | `notes/kp/01-变量与常量（作用域、_、iota）.md` | `go-learning/cmd/kp/vars-const-scope-iota` |
| KP 02 | 基础类型（转换、格式化、表达式） | `notes/kp/02-基础类型（转换、格式化、表达式）.md` | `go-learning/cmd/kp/basic-types` |
| KP 03 | 字符串基本操作（转义、格式化、Builder、比较、常用方法） | `notes/kp/03-字符串基本操作（转义、格式化、Builder、比较、常用方法）.md` | `go-learning/cmd/kp/string-basics` |
| 索引 | 当前进度 / 新对话恢复 / 记忆策略 | `notes/progress.md`、`notes/context-pack.md`、`notes/memory-strategy.md` | 无 |
| 原子记忆 | 术语 / 模式 / 坑点 | `notes/glossary.md`、`notes/patterns.md`、`notes/pitfalls.md` | 无 |

## 下一步（建议）
- 如果继续系统化主线：从 `notes/day06-gin-engineering.md` 重新起一条新的 Day 笔记
- 如果继续点播补课：沿用 `notes/kp/NN-*.md` 持续追加
