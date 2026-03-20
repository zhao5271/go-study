# knowledge-point-notes defaults

## 默认工作区与目录
- 默认工作区：`/Users/zhang/Desktop/go-study/codex`
- 知识点笔记目录：`/Users/zhang/Desktop/go-study/codex/notes/kp`
- 可运行示例目录：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/kp`
- 工作区记忆方案：`/Users/zhang/Desktop/go-study/codex/notes/memory-strategy.md`
- Obsidian 属性规范：`/Users/zhang/Desktop/go-study/codex/notes/Obsidian-属性规范.md`
- 上下文压缩包：`/Users/zhang/Desktop/go-study/codex/notes/context-pack.md`
- 外部记忆索引：`/Users/zhang/Desktop/go-study/codex/notes/progress.md`
- 一致性约束：`/Users/zhang/Desktop/go-study/codex/notes/glossary.md`、`/Users/zhang/Desktop/go-study/codex/notes/patterns.md`、`/Users/zhang/Desktop/go-study/codex/notes/pitfalls.md`

## 记忆方案入口
- 若工作区存在 `notes/memory-strategy.md`，以它为准。
- `references/defaults.md` 只负责默认路径、命名与 NotebookLM 配置；不要把完整记忆策略继续堆在这里。

## Obsidian 集成入口
- 若工作区存在 `notes/Obsidian-属性规范.md`，新建/更新笔记时按其 frontmatter 规范输出。
- 默认新知识点笔记至少带：`type`、`domain`、`topic`、`topic_zh`、`stage`、`status`、`review_cycle`、`tags`。

## 操作识别
- 新建/默认更新：`知识点：...`、`做知识点笔记：...`、`点播知识点：...`、`补一节知识点笔记：...`
- 明确更新：`更新知识点笔记：...`、`更新 xxx`、`补充 xxx`
- 列表：`列出知识点笔记`
- 只同步记忆：`同步外部记忆`

## 文件命名规则
- 笔记文件：`NN-<简短中文>.md`
- `NN` 使用两位数字，取 `notes/kp/*.md` 里现有最大编号 + 1
- 若更新已有知识点，优先复用现有文件，不新开编号
- 标题要求：短、可搜索、尽量避免重复词
- demo 目录名使用英文 `kebab-case`，例如 `context-cancel`、`slice-copy`、`gin-middleware`

## 笔记推荐结构
1. `# 标题`
2. `## TL;DR（可放入 progress/context-pack）`
3. `## 关键词`
4. `## 定义`
5. `## 为什么重要`
6. `## 边界 / 取舍 / 常见坑`
7. `## 工程落地`
8. `## 关键代码`
9. `## 知识点运用示例`
10. `## 关联复习`（仅在需要复述前置知识时添加）
11. `## References`（仅 web fallback 用到时添加）

## 代码与输出规则
- 只放关键代码，不贴整文件
- `fmt.Print*`、`log.Print*`、`panic`、命令行输出、HTTP 响应、`curl` 示例都要在紧邻位置写典型输出注释
- 输出不稳定时，必须标注“输出可能变化/不固定”，并说明原因
- 练习题与参考答案必须就地放在对应知识点下面，不集中放到文末

## NotebookLM 默认配置
- NotebookLM skill 目录：`/Users/zhang/.cc-switch/skills/notebooklm`
- 默认 Notebook URL：`https://notebooklm.google.com/notebook/1e4b57b8-8e53-4fbe-a322-a4dfd1e2725d`
- 必须通过 `run.py` 调用：

```bash
cd /Users/zhang/.cc-switch/skills/notebooklm

# 1) 检查认证状态
python scripts/run.py auth_manager.py status

# 2) 首问：覆盖定义、why、坑点、工程落地、示例
python scripts/run.py ask_question.py \
  --notebook-url "https://notebooklm.google.com/notebook/1e4b57b8-8e53-4fbe-a322-a4dfd1e2725d" \
  --question "<知识点>：请给出定义、为什么重要、边界与取舍、常见坑、后台管理 API 的工程落地建议，并给最小可运行示例（越小越好）"

# 3) 缺项追问：只补缺的
python scripts/run.py ask_question.py \
  --notebook-url "https://notebooklm.google.com/notebook/1e4b57b8-8e53-4fbe-a322-a4dfd1e2725d" \
  --question "针对上次回答缺失的部分：<列出缺失项>，请只补这些缺失项，并给更贴近后台管理 API 的例子"
```

## 缺项清单
- 定义是否清楚
- 为什么重要是否能落到真实 API 项目
- 是否覆盖边界、取舍、常见误解
- 是否有工程落地建议
- 是否有最小可运行示例
- 示例是否标了输出注释

## Web fallback 规则
- NotebookLM 不可用、资料不完整、或没有 runnable / pitfalls / 工程落地时，直接补 web
- 官方优先级：Go Spec → `go.dev` → `pkg.go.dev`
- 社区补充只用于解释、经验、最佳实践，不替代官方定义
- 用过 web fallback 后，在 `## References` 里写：
  - 链接
  - 一句话用途
  - 标注“官方”或“社区”
  - 如果 NotebookLM 失败，再补一句：`NotebookLM 查询失败（原因：...），本节结论以官方资料补齐。`

## 默认回复形态（省 token）
- 文件清单
- 5–10 条 TL;DR
- 1–3 段关键代码或命令
- 如何运行
- References 说明（若有）
- 外部记忆 delta
- 用户明确要求“完整笔记”前，不粘贴整篇内容
