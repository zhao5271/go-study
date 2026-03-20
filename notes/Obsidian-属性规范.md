---
type: reference
domain: go
role: obsidian-schema
status: active
tags:
  - go
  - obsidian
  - schema
created: 2026-03-20
updated: 2026-03-20
---

# Obsidian 属性规范

## 目标
- 让知识点笔记、索引、复习日志在 Obsidian 里能稳定筛选、分组、做 Bases 视图。
- 字段尽量少而稳，避免同义字段越积越多。

## 推荐字段

| 字段 | 用途 | 推荐值 |
|---|---|---|
| `type` | 笔记类型 | `kp` / `lesson` / `index` / `reference` / `dashboard` / `journal` |
| `domain` | 知识域 | 固定 `go` |
| `topic` | 英文 slug | 如 `string-basics`、`context-timeout` |
| `topic_zh` | 中文主题 | 如 `字符串基本操作` |
| `stage` | 学习阶段 | `foundation` / `http` / `db` / `engineering` / `deploy` / `review` |
| `status` | 当前状态 | `active` / `evergreen` / `archived` |
| `review_cycle` | 复习频率 | `daily` / `weekly` / `biweekly` / `monthly` |
| `source` | 资料来源 | `notebooklm` / `official` / `community` / `practice` |
| `tags` | Obsidian 标签 | 例如 `go`、`kp`、`foundation` |
| `created` | 创建日期 | `YYYY-MM-DD` |
| `updated` | 更新日期 | `YYYY-MM-DD` |

## 使用规则
- `type`、`domain`、`status` 必填。
- `kp` 笔记额外要求：`topic`、`topic_zh`、`stage`、`review_cycle`。
- `index/reference/dashboard` 笔记不要求 `topic`。
- 字段命名保持稳定：不要混用 `kind/type`、`area/domain`、`review/review_cycle`。

## 建议映射
- `kp/*.md` → `type: kp`
- `day*.md` → `type: lesson`
- `progress/context-pack` → `type: index`
- `glossary/patterns/pitfalls/memory-strategy` → `type: reference`
- `知识库首页.md` → `type: dashboard`
- `daily/*.md` → `type: journal`
