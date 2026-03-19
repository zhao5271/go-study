---
name: knowledge-point-notes
description: Create detailed “知识点笔记” for Go learning (notes/go/kp) with optional runnable demos (go-learning/cmd/kp), and keep external memory + context-pack updated for token compression. Triggers when user says “知识点：…/做知识点笔记/点播知识点/补一节知识点笔记”.
---

# knowledge-point-notes

把用户给的“知识点”沉淀成 **可长期复用的知识库笔记**，并同步维护 **外部记忆索引**，用于解决对话 token 膨胀问题。

## 重开窗口/新对话启动（必须做）
1) 读上下文压缩包：`notes/go/context-pack.md`
2) 读外部记忆索引：`notes/go/progress.md`
3) 读一致性约束：`notes/go/glossary.md`、`notes/go/patterns.md`、`notes/go/pitfalls.md`

## 工作区与目录约定（避免迷路）
- 仓库根：`/Users/zhang/Desktop/go-study/codex`
- 点播笔记目录：`notes/go/kp/`
  - 文件名：`NN-<简短中文>.md`（NN 两位数字，递增；标题尽量短、便于搜索）
- 可运行示例目录（可选）：`go-learning/cmd/kp/<topic-slug>/main.go`
  - `<topic-slug>` 用小写英文 `kebab-case`（例如 `basic-types` / `vars-const-scope-iota`）
  - 运行都在 Go module 内执行：先 `cd go-learning`

## 使用场景（When to use）
- 用户说“知识点：xxx / 点播知识点 / 做知识点笔记 / 补一节知识点笔记”
- 你要把零散对话内容沉淀为文件知识库（可检索、可复习、可复用）

## 输出与写作硬规则（必须遵守）
### 1) 笔记要“自包含”
- **不要在笔记里出现** `/Users/.../xxx.md:行号` 这类引用路径（读者不应该被迫跳出去看）。
- 需要关联旧知识时：在本笔记加 `## 关联复习`，用 3–12 行把前置知识“复述到位”（允许提到文件名，但不要引用行号/让读者跳转）。

### 2) 笔记要“详细，但不啰嗦”
- 目标是“长期知识库”，不是提纲：要写清楚**边界、取舍、常见误解、工程落点**。
- 讲解步骤不必教条化，但至少覆盖：定义/为什么重要/关键坑与取舍/工程落地/自检或练习。

### 3) 代码片段规则（笔记里的代码）
- **只放主要代码**（与知识点强相关），不要贴整文件，更不要塞无关分隔输出（例如 `fmt.Println("----")`）。
- 代码里只要会产生输出（fmt/log/panic/HTTP 响应/curl 示例等），必须在同一行或紧邻位置注释**典型输出**。
  - 若输出不确定：必须注释“输出可能变化/不固定 + 原因”（例如平台位宽、map 顺序、错误信息差异等）。
- “练习/自检”不要集中到文末：统一使用小标题 **“知识点运用示例”** 放在对应知识点下面，并紧跟参考答案。

### 4) 可运行示例规则（go-learning/cmd/kp，可选但默认优先）
- 能用标准库做的就只用标准库。
- 示例要**最小可运行**，不引入“还没学过的大概念”；确实需要前置时：
  - 先用 1–3 句解释前置作用
  - 或者明确说明“本节不做可运行示例”的原因（避免硬塞导致更难理解）

## 资料来源策略（NotebookLM first, fallback web）
1) 先查 NotebookLM（参考资料库；能覆盖就优先遵从其材料与结构）
2) 若缺关键项（why / runnable / pitfalls / 工程落地）或 NotebookLM 失败：用 web fallback 补齐
   - 关键结论以官方为准：Go Spec / go.dev / pkg.go.dev
3) 只要用了 web fallback：笔记末尾追加 `## References`（链接 + 一句话用途 + 官方/社区标记）

### NotebookLM 固定执行步骤（把它当“本地知识库查询”）
Notebook URL（固定）：
- `https://notebooklm.google.com/notebook/1e4b57b8-8e53-4fbe-a322-a4dfd1e2725d`

命令（必须通过 `run.py` 包装器执行）：
```bash
cd /Users/zhang/.cc-switch/skills/notebooklm

# 1) 检查登录状态
python scripts/run.py auth_manager.py status

# 2) 如未登录：做一次可见浏览器登录（用户手动登录 Google）
python scripts/run.py auth_manager.py setup

# 3) 首问：围绕“知识点”要齐全（定义/why/坑/工程落地/示例）
python scripts/run.py ask_question.py \
  --notebook-url "https://notebooklm.google.com/notebook/1e4b57b8-8e53-4fbe-a322-a4dfd1e2725d" \
  --question "<知识点>：请给出定义、为什么重要、常见坑（含边界/取舍）、工程落地建议，并给最小可运行示例（越小越好）"

# 4) follow-up loop：如果缺项，继续追问“只补缺失项”，直到足够写笔记
python scripts/run.py ask_question.py \
  --notebook-url "https://notebooklm.google.com/notebook/1e4b57b8-8e53-4fbe-a322-a4dfd1e2725d" \
  --question "针对上次回答缺失的部分：<列出缺失项>，请补齐，并给出更贴近后台管理 API 的例子"
```

NotebookLM 失败兜底（超时/未登录/打不开/连接中断）：
- 不要卡住，直接进入 web fallback（官方优先）补“缺失项”
- 在笔记 `## References` 里加一句：`NotebookLM 查询失败（原因：...），本节结论以官方资料补齐。`

## 外部记忆更新（强制）
**只要生成/更新了 `notes/go/kp/*.md` 或 `go-learning/cmd/kp/**`，就必须同步更新外部记忆：**
- `notes/go/progress.md`：更新“点播笔记”列表（新增/更新哪一篇）
- `notes/go/glossary.md`：新增术语（每条 3–8 行）
- `notes/go/patterns.md`：沉淀可复用模板（参数解析/错误包装/分页/权限位等）
- `notes/go/pitfalls.md`：记录坑点与规避（1–3 行/条）
- `notes/go/context-pack.md`：更新“当前进度/最近更新/如何继续”，让新对话只贴这一份就能恢复上下文
- 同步写入外部记忆图谱（MCP memory）：为新笔记创建/更新节点，并把本次“新增术语/模式/坑点”写成 observation

若本次没有新增术语/模式/坑点，也要在回复里明确写一句“外部记忆已检查，无需新增”。

## 上下文压缩协议（重点优化）
目标：避免对话超过 258k token，且重开对话能无缝继续。

1) **对话里默认不粘贴整篇笔记**
   - 默认只输出：文件清单 + TL;DR + 关键代码片段 + 外部记忆 delta
   - 用户说“输出完整笔记”时，才把整篇笔记内容贴出来
2) **把“可复用信息”强制落盘**
   - 新术语 → `glossary.md`
   - 新套路 → `patterns.md`
   - 新坑 → `pitfalls.md`
   - 本次进度/入口 → `progress.md` + `context-pack.md`
3) **每篇 kp 笔记都加“压缩锚点”**
   - 建议在笔记靠前位置加 `## TL;DR（可放入 progress/context-pack）`（5–10 条）
   - 建议加 `## 关键词`（便于搜索/检索）

## 操作面板（可视化，可选）
当用户没说清楚具体要干什么时，展示并让他选 1 个；用户明确就直接执行。

| 操作 | 触发方式 | 输入 | 输出 |
|---|---|---|---|
| A. 新建知识点笔记（默认） | “知识点：xxx / 做知识点笔记：xxx” | 知识点名称 | `notes/go/kp/NN-*.md` +（可选）`go-learning/cmd/kp/<slug>/main.go` |
| B. 更新已有知识点笔记 | “更新知识点笔记 NN / 更新 xxx” | NN 或文件名 | 更新对应笔记 + 外部记忆同步 |
| C. 查看知识点列表 | “列出知识点笔记” | 无 | 列出 `notes/go/kp` 文件名 |
| D. 只更新外部记忆 | “同步外部记忆” | 无 | 校验并更新 progress/glossary/patterns/pitfalls/context-pack + memory 图谱 |
