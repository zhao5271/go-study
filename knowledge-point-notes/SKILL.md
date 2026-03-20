---
name: knowledge-point-notes
description: Create or update Go knowledge-point notes with NotebookLM-first research, web fallback, optional runnable demos, and external-memory sync. Use when the user asks “知识点：…”, “做知识点笔记”, “更新知识点笔记”, “列出知识点笔记”, or “同步外部记忆”.
---

# knowledge-point-notes

把用户给的 Go 知识点沉淀成**可长期复用的笔记 + 外部记忆**，默认优先写文件，再用简短摘要回话，避免 token 膨胀。

## 先读默认约定
- 开始前先读 `references/defaults.md`。
- 如果目标工作区里存在 `notes/memory-strategy.md`，继续读取它，并把它视为该工作区的记忆方案权威说明。
- 如果目标工作区里存在 `notes/Obsidian-属性规范.md`，继续读取它，并在新建/更新笔记时遵守其中的 frontmatter 规范。
- 如果用户在请求里明确给了工作区、笔记目录、Notebook URL 或输出方式，按用户显式要求覆盖默认值。
- 这个 skill 是全局可用的，但默认落到当前 Go 学习仓库。

## 支持的操作
- `知识点：<topic>` / `做知识点笔记：<topic>`：为单个知识点新建或更新笔记。
- `更新知识点笔记：<topic|编号>`：补充已有笔记，保留原编号与已有有效内容。
- `列出知识点笔记`：列出当前知识点笔记，结果要简短、可搜索。
- `同步外部记忆`：只校验并同步 `progress/context-pack/glossary/patterns/pitfalls` 与 MCP memory。
- 默认一条请求只做一个知识点；只有用户明确说“专题 / 汇总 / 多个知识点一起整理”时，才输出专题笔记。

## 工作流
1. 先读取当前上下文：
   - `notes/memory-strategy.md`（如果存在）
   - `notes/Obsidian-属性规范.md`（如果存在）
   - `notes/context-pack.md`
   - `notes/progress.md`
   - 若是更新操作，再读对应旧笔记和 demo
   - 只有当主题确实需要术语/模式/坑点一致性时，再按需读取 `glossary.md`、`patterns.md`、`pitfalls.md`
2. 解析目标路径：
   - 新笔记使用下一个 `NN-<简短中文>.md`
   - 更新已有知识点时，优先复用已有编号/文件名，不重复编号
   - 可运行示例按需写到 `go-learning/cmd/kp/<topic-slug>/main.go`
3. NotebookLM first：
   - 复用现有 NotebookLM skill 的 `run.py` 包装器，不要复制脚本
   - 先检查认证，再围绕“定义 / 为什么重要 / 边界与取舍 / 常见坑 / 工程落地 / 最小示例”提问
   - 首答缺项时，立刻用 follow-up 只追问缺失项，不要直接交半成品
4. Web fallback：
   - NotebookLM 超时、未登录、连接失败，或内容仍不完整时，改用 web 补齐
   - 关键结论优先采用 Go 官方资料，其次再补高质量社区资料
   - 只要用了 web fallback，就在笔记末尾追加 `## References`
5. 写笔记：
   - 笔记必须自包含，不能写 `/Users/.../file:line` 这类跳转引用
   - 若工作区配置了 Obsidian 属性规范，先写 frontmatter，再写正文；新字段优先复用既有值，避免同一库里命名漂移
   - 靠前加入 `## TL;DR（可放入 progress/context-pack）` 和 `## 关键词`
   - 重点写清楚定义、为什么重要、边界/取舍、常见误解、工程落地
   - 代码片段只保留关键部分；凡是会产生输出的代码/命令/HTTP 响应，都要在附近写典型输出注释；如果输出不固定，要说明原因
   - 练习题不要堆到文末，统一在对应知识点下用 `## 知识点运用示例`，并紧跟参考答案
6. 写 demo：
   - 优先标准库，保持最小可运行
   - 只有 demo 真能帮助理解时才创建/更新
   - 如果必须引入太多前置概念，就明确说明“本节不提供可运行示例”和原因
7. 同步记忆：
   - 优先遵循工作区的 `notes/memory-strategy.md`
   - 若无工作区方案，则默认：笔记或 demo 有变更时更新 `progress/context-pack`；`glossary/patterns/pitfalls` 与 MCP memory 只在存在高价值 delta 时更新
   - 若本次没有新增术语/模式/坑点，也要在回复里明确写“外部记忆已检查，无需新增”

## 回复约定
- 默认只回：文件清单、TL;DR、关键代码片段、运行命令、References 说明、memory delta。
- 用户明确说“输出完整笔记”时，才粘贴整篇内容。
- `列出知识点笔记` 与 `同步外部记忆` 只返回最小必要结果，不展开成长文。
