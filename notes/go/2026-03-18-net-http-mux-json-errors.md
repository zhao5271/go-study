# 2026-03-18：net/http ServeMux + 统一 JSON 错误返回

## 今日目标
- 用 `net/http` 写一个能“交付给前端联调”的最小后台 API：健康检查 + 用户列表
- 把 **404/405/参数错误** 统一成 JSON，前端好处理、后端好排障

## 背景对照（只说会踩坑的）
- 在框架（如 Express/Nest）里，“路由=方法+路径”是默认能力；而原生 `net/http` 里 **ServeMux 默认只看路径**，方法校验要你自己做（或 Go1.22+ 用新 pattern 语法）。

## 关键结论
- `http.HandlerFunc` 本质是把 `func(w, r)` 适配成 `http.Handler`
- `http.NewServeMux()` + 显式注册路由，利于测试与隔离（不要滥用全局 `http.HandleFunc`）
- 统一 JSON 返回：把“写 Header/Status/Body”收敛到一个 `writeJSON`，避免每个 handler 复制粘贴
- 404/405/400 等非业务错误也要“稳定格式”，否则前端会出现大量分支处理

## 代码清单
- `go-learning/cmd/day04_02_http_admin_api_min/main.go`

## 常见坑
- `Content-Type` 必须在 `WriteHeader/Write/Encode` 前设置
- `http.Error(...)` 后必须 `return`，否则后续继续写响应会导致“header already written”等问题
- ServeMux 的末尾 `/` 会影响匹配（精确 vs 前缀），上线前要用 `curl` 验证

## 面试问法（自测）
1) `http.Handler` 和 `http.HandlerFunc` 有什么关系？为什么要用它？
2) 你会怎么做统一错误返回？为什么不在每个 handler 里手写？
3) 为什么要给 API 做 `/api/v1/health`？它在 Docker/K8s 里有什么用？

## 下一步
- 把 `/api/v1/users` 接上 MySQL：分页 + 模糊搜索 + 统一错误码
- 引入 `context.WithTimeout`：避免慢查询/下游卡死拖垮服务

