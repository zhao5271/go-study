---
type: reference
domain: go
role: patterns
status: active
tags:
  - go
  - reference
  - patterns
created: 2026-03-20
updated: 2026-03-20
---

# Go 工程模式库（外部记忆）

> 目的：把“可复用的工程套路”写下来，避免每次都从头推导。

## 错误处理（服务/仓库层）
- 函数签名优先返回 `(value, error)`，在边界处统一做错误映射。
- 常见模式：下层返回原始错误，上层按是否需要 `errors.Is/As` 来决定用 `%w` 还是 `%v`。

## HTTP Handler（提前约定）
- handler 只负责：解析输入 → 调 service → 把 error 映射成 HTTP 响应。
- 内部逻辑不要直接写响应；统一返回 error，在入口层转换。

## 统一 JSON 响应（最小版）
- 约定统一响应结构：`{code, message, data}`
- 统一 `writeJSON` 只做 Content-Type + status + encode。
- 统一 `writeError` 做错误码与 HTTP status 分离，便于前端稳定处理。

## 分页 query（最小版）
- `page>=1`；`1<=size<=100`（统一边界）
- 用统一的分页解析 helper 处理 `page/size`，失败就返回 400 + `INVALID_QUERY`

## DB 访问（最小版）
- DSN 从环境变量读取，避免把地址、账号、参数写死在代码里。
- 所有 DB 操作都带 `context.WithTimeout`；列表查询优先用 `COUNT(*)` + `LIMIT/OFFSET` 先跑通，再谈优化。

## Repo/Service 分层（最小版）
- service 持有接口并通过构造函数注入，repo 接口表达业务语义，不暴露 SQL/驱动细节。
- 横切能力（log/metrics）优先用装饰器：`type LoggingRepo struct { Repo }`。

## 编译期接口校验（防止漏实现）
- 场景：重构 repo/service、拆包或改接收者类型后，避免“以为实现了接口，实际没实现”。
- 模式：用 `var _ IFace = (*Impl)(nil)` 做编译期校验，把失败前移到编译阶段。

## Table-driven tests 模板
- 用例表统一包含：名字、输入、期望输出、期望错误语义。
- 循环里统一 `t.Run(tt.name, ...)`，把断言结构保持一致，减少漏测边界值。

## 枚举常量（iota）
- 用 `type` + `const` 表达稳定枚举（角色/状态/错误类别），避免 magic number/string。
- 如果 `0` 代表“未设置/非法”，通常从 `iota + 1` 开始更稳。

## 参数解析（默认值 + 转换 + 边界）
- 场景：分页 `page/size`、筛选 `status`、导出 `limit` 等。
- 模式：先给默认值，再做 `Atoi/ParseInt`，最后统一做范围校验；失败后交给上层统一映射错误码/HTTP status。

## 安全数值转换（int64 → int）
- 场景：DB 扫描/协议字段是 `int64`，但业务里需要 `int`（下标/长度/分页）。
- 模式：在边界处显式转换，并在 32 位平台上先做溢出判断。

## 位标志（权限/状态）
- 场景：用户权限（READ/WRITE/EXPORT）、状态位（禁用/冻结/需要改密）等。
- 模式：用自定义整数类型保存 mask，用 `1 << iota` 定义 flag，用 `(mask & flag) != 0` 判断。

## 字符串拼接（`strings.Builder`）
- 场景：审计日志一行文本、导出 CSV、构造较长错误信息（多段拼接/循环拼接）。
- 模式：预估长度就 `Grow`，循环内用 `WriteString/WriteByte`，最后统一 `String()`。

## 动态 WHERE（只拼结构，值走参数化）
- 场景：列表分页检索（status/keyword/role 等可选过滤）。
- 模式：SQL 里只拼接 `AND ... = ?` 这类结构，把值放进 `args`，交给 driver 参数化，避免注入。
