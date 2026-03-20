---
type: reference
domain: go
role: memory-strategy
status: active
tags:
  - go
  - reference
  - memory
created: 2026-03-20
updated: 2026-03-20
---

# Go 学习外部记忆方案

> 目标：把“恢复上下文”和“保存知识”拆开，减少重复写入、减少重复读取、减少 token 浪费。

## 一、记忆分层

### 1) 详细内容层：`notes/day*.md` + `notes/kp/*.md`
- 这是**唯一的详细知识源**。
- 真正完整的解释、代码、练习、参考资料只放在这里。
- 其他记忆文件只保留索引、摘要、术语、模式、坑点，不重复抄整段内容。

### 2) 进度索引层：`notes/progress.md`
- 这是**当前学习进度与产物目录的权威索引**。
- 记录当前进度、已完成主题、知识点笔记清单、下一步建议。
- 适合回答“学到哪了”“有哪些笔记”“下一步做什么”。

### 3) 恢复上下文层：`notes/context-pack.md`
- 这是**开新对话的最小恢复包**。
- 只保留：当前路线、当前进度、最近点播、下一步建议、关键文件入口。
- 不再重复大段稳定规则；稳定规则应放在仓库 `AGENTS.md` 与 skill 里。

### 4) 可复用原子层：`notes/glossary.md`、`notes/patterns.md`、`notes/pitfalls.md`
- 这里只存**可跨主题复用**的原子信息。
- `glossary.md`：术语定义。
- `patterns.md`：可复用工程套路。
- `pitfalls.md`：高复用踩坑提醒。
- 单篇笔记独有、不会复用的内容，不要硬塞进这三份文件。

### 5) 稀疏远端索引层：MCP memory
- 只存**高价值 delta**，用于跨对话快速检索。
- 适合存：新增笔记标题、超短摘要、可复用术语/模式/坑点。
- 不要把整篇笔记、长摘要、完整代码再次镜像进去。

## 二、读取策略

### 新建/更新知识点笔记
- 先读：`notes/memory-strategy.md`、`notes/context-pack.md`、`notes/progress.md`
- 若是更新已有笔记，再读目标笔记与对应 demo
- 只有当主题涉及既有术语/模式/坑点复用时，再按需读 `glossary.md`、`patterns.md`、`pitfalls.md`
- 不要为了一个知识点默认把所有记忆文件全量读一遍

### 列出知识点笔记
- 只读 `notes/progress.md`
- 除非索引明显过期，否则不要读取其他记忆文件

### 同步外部记忆
- 先读 `notes/progress.md`、`notes/context-pack.md`
- 再只读取本次可能需要变更的原子层文件
- 只有发现索引漂移或命名冲突时，才回读具体笔记

## 三、写入策略

### 每次笔记或 demo 发生实质变更时，必更
- `notes/progress.md`
- `notes/context-pack.md`

### 只有出现“可复用 delta”时才更新
- `notes/glossary.md`：这个术语未来大概率还会复用
- `notes/patterns.md`：这个套路未来大概率能复用到别的主题或项目
- `notes/pitfalls.md`：这个坑具备普遍性，不只是当前例子里的细节
- MCP memory：这次新增了值得跨对话记住的高价值信息

### 不必触发记忆更新的情况
- 只改了措辞、排版、标题格式
- 只是把已有内容重写得更顺，但没有新增知识 delta
- 当前内容只服务单篇笔记，不具备复用价值

## 四、低价值/应删减的部分

### 1) `context-pack.md` 里重复稳定规则
- 当前 `context-pack.md` 里曾放了角色、教学规则、资料策略等稳定说明。
- 这些内容已经被仓库 `AGENTS.md`、skill、固定工作流覆盖，再重复一次会浪费 token。
- 结论：`context-pack.md` 只保留“恢复会话所需的最小状态”。

### 2) 每次点播都强制更新五层记忆
- 如果每次都改 `progress/context-pack/glossary/patterns/pitfalls/MCP memory`，维护成本高、噪音大、漂移概率也高。
- 结论：改成“索引层必更，原子层按 delta 更新”。

### 3) 同一段摘要在多处复制
- 把同一段 TL;DR 同时写进笔记、`progress.md`、`context-pack.md`、MCP memory，收益很低。
- 结论：详细摘要留在笔记；索引与远端只放更短的摘要或指针。

### 4) `列出知识点笔记` 读取过多文件
- 这个操作本质是“列目录 + 简短描述”，只读 `progress.md` 足够。
- 结论：不需要顺手把 glossary/patterns/pitfalls 也读进来。

### 5) 为了“完整”而强行创建 demo
- 有些主题 demo 价值很高；有些主题只会引入更多前置知识，反而让笔记更重。
- 结论：demo 只在确实能降低理解成本时创建。

## 五、维护准则

### `progress.md` 应该写什么
- 当前学习阶段
- 已完成主题与产物入口
- 知识点笔记清单
- 下一步建议

### `context-pack.md` 应该写什么
- 最近进度
- 最近点播
- 下一步从哪里继续
- 关键索引文件入口

### `glossary/patterns/pitfalls` 不该写什么
- 整段教学解释
- 单次实验细节
- 只在当前一篇笔记里出现、不会复用的描述

## 六、执行建议
- 默认让 `knowledge-point-notes` 先遵守这份文件，再做具体写入。
- 如果后续你新增新的知识库目录或改文件命名规则，优先改这份文件，再改 skill 默认值。
