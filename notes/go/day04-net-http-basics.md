# Day 04.1：net/http 入门（ServeMux/Handler）+ JSON 响应与错误码（为后台管理 API 打底）

> 默认贯穿项目：后台管理 API（RBAC + 登录鉴权 + 列表分页检索 + CRUD + 审计日志 + Docker 部署）

## 你更偏向哪种学习目标？（不回答按默认继续）
1) 作品集项目交付（默认）
2) 面试冲刺
3) 基础深挖

## 你是否能使用 Docker？（不回答按默认继续）
- 默认：能

## 你 MySQL 基础？（不回答按默认继续）
- 熟悉 / 一般（默认）/ 没有

---

## 知识点 1：`http.Handler` + `http.ServeMux`（最小可交付的 HTTP 服务骨架）

### B. 一句话定义
`Handler` 是“处理一次 HTTP 请求的函数/对象”，`ServeMux` 是“把 URL 路径分发到不同 Handler 的路由表”。

### C. 为什么重要（不做会怎样）
后台管理 API 的所有功能（登录、列表、增删改）最终都要落在 HTTP 路由与 Handler 上；如果骨架不清晰，很快会变成“每个接口各写各的，日志/错误/超时/鉴权全散落”，后期维护成本爆炸。

### D. 重难点拆解（2–4 条）
1) **方法语义**：同一路径不同 Method（GET/POST/PATCH/DELETE）语义不同，别用“动作型 URL”替代（如 `/createUser`）。
2) **返回 405**：不支持的 Method 要明确返回 405（否则前端/调用方会误判为 404 或 500）。
3) **可观测性最小集**：至少要有启动日志 + 请求日志的落点（后面再加 request id、耗时、结构化日志）。

### E. 业务场景落地（后台管理 API）
- `/api/v1/health`：容器探活与发布健康检查（K8s/Docker Compose 都需要）。

### F. 代码示例（最小可运行）

运行：
- `cd /Users/zhang/Desktop/go-study/codex/go-learning && PORT=18080 go run ./cmd/day04_01a_http_mux`

验证（另开一个终端）：
- `curl -i http://localhost:18080/health`
  - 典型响应：
    - `HTTP/1.1 200 OK`
    - body：`ok`（文本）
- `curl -i -X POST http://localhost:18080/health`
  - 典型响应：
    - `HTTP/1.1 405 Method Not Allowed`
    - body：`method not allowed`

代码（全文）：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day04_01a_http_mux/main.go`
```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// 健康检查：给部署/容器编排做 readiness/liveness 用
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed) // Output: HTTP 405 "method not allowed\n"
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintln(w, "ok") // Output: ok
	})

	addr := ":8080"
	log.Printf("listening on %s", addr) // Output: 2006/01/02 15:04:05 listening on :8080 (输出可能变化/不固定：包含时间戳)
	err := http.ListenAndServe(addr, mux)
	log.Printf("server stopped: %v", err) // Output: 2006/01/02 15:04:05 server stopped: listen tcp :8080: bind: address already in use (输出可能变化/不固定：取决于环境与错误)
}
```

### G. 怎么运行（命令 + 预期现象）
- 启动：`cd go-learning && go run ./cmd/day04_01a_http_mux`
- 预期现象：终端打印 listening；curl `/health` 返回 200 + `ok`

### H. 练习题（1–3 题，覆盖边界条件）

#### 练习 1：把 `/health` 改成版本化路径 `/api/v1/health`
- 验收标准：
  - `curl -i http://localhost:8080/api/v1/health` 返回 200
  - `curl -i http://localhost:8080/health` 返回 404（或被你重定向，但要解释）

### I. 参考答案（可运行）
直接看下一个知识点的示例（知识点 2），它已经使用了 `/api/v1/health`。

---

## 知识点 2：统一 JSON 响应 + 错误码（让前端/排障/监控都省心）

### B. 一句话定义
把所有接口的成功/失败都包装成统一 JSON 结构（含 `code/message/data`），并用正确的 HTTP Status Code 表达“请求在 HTTP 语义上成功还是失败”。

### C. 为什么重要（不做会怎样）
- 前端没法稳定处理：有时是文本、有时是 JSON、字段还不一致，会导致大量 if/else 与隐藏 bug。
- 排障困难：监控/日志里很难聚合相同错误；你也很难做“错误码 → 告警/埋点/灰度策略”。

### D. 重难点拆解（2–4 条）
1) **HTTP status vs 业务 code**：HTTP status 表达协议层成功/失败；业务 code 表达业务原因（如参数非法、鉴权失败）。
2) **参数校验与边界**：分页 `page/size` 的范围要有限制（比如 size <= 100），避免单请求拖垮服务。
3) **一致性**：405/400/404 都要有一致 JSON 结构，否则调用方仍要写分支。

### E. 业务场景落地（后台管理 API）
- 列表分页：`GET /api/v1/users?page=1&size=2`（后面会扩展到 `search/role` 等过滤）

### F. 代码示例（最小可运行）

运行：
- `cd /Users/zhang/Desktop/go-study/codex/go-learning && PORT=18080 go run ./cmd/day04_01b_json_errors`

验证（另开一个终端）：
- `curl -s http://localhost:18080/api/v1/health`
  - 典型响应：
    - `{"code":0,"message":"OK","data":{"ok":true}}`
- `curl -s "http://localhost:18080/api/v1/users?page=1&size=2"`
  - 典型响应（示例）：
    - `{"code":0,"message":"OK","data":{"items":[{"id":1,"name":"Alice","role":"admin"},{"id":2,"name":"Bob","role":"viewer"}],"page":1,"size":2,"total":4}}`
- `curl -s "http://localhost:18080/api/v1/users?page=0"`
  - 典型响应：
    - `{"code":10002,"message":"INVALID_QUERY"}`

代码（全文）：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day04_01b_json_errors/main.go`
```go
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type ListUsersData struct {
	Items []User `json:"items"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Total int    `json:"total"`
}

var errBadQuery = errors.New("bad query")

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v) // Output: {"code":0,"message":"OK","data":{...}}\n
}

func writeError(w http.ResponseWriter, status int, code int, message string) {
	writeJSON(w, status, APIResponse{Code: code, Message: message})
}

func parsePageSize(r *http.Request) (page int, size int, err error) {
	page = 1
	size = 20

	q := r.URL.Query()
	if raw := strings.TrimSpace(q.Get("page")); raw != "" {
		n, convErr := strconv.Atoi(raw)
		if convErr != nil || n < 1 {
			return 0, 0, errBadQuery
		}
		page = n
	}
	if raw := strings.TrimSpace(q.Get("size")); raw != "" {
		n, convErr := strconv.Atoi(raw)
		if convErr != nil || n < 1 || n > 100 {
			return 0, 0, errBadQuery
		}
		size = n
	}
	return page, size, nil
}

func main() {
	users := []User{
		{ID: 1, Name: "Alice", Role: "admin"},
		{ID: 2, Name: "Bob", Role: "viewer"},
		{ID: 3, Name: "Carol", Role: "editor"},
		{ID: 4, Name: "Dave", Role: "viewer"},
	}

	mux := http.NewServeMux()

	// 真实项目里建议统一加 /api/v1 做版本前缀，避免未来破坏性变更无处安放。
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, 10001, "METHOD_NOT_ALLOWED") // Output: {"code":10001,"message":"METHOD_NOT_ALLOWED"}\n
			return
		}
		writeJSON(w, http.StatusOK, APIResponse{Code: 0, Message: "OK", Data: map[string]bool{"ok": true}}) // Output: {"code":0,"message":"OK","data":{"ok":true}}\n
	})

	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, 10001, "METHOD_NOT_ALLOWED") // Output: {"code":10001,"message":"METHOD_NOT_ALLOWED"}\n
			return
		}

		page, size, err := parsePageSize(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, 10002, "INVALID_QUERY") // Output: {"code":10002,"message":"INVALID_QUERY"}\n
			return
		}

		total := len(users)
		start := (page - 1) * size
		if start > total {
			start = total
		}
		end := start + size
		if end > total {
			end = total
		}

		log.Printf("list_users page=%d size=%d start=%d end=%d", page, size, start, end) // Output: 2006/01/02 15:04:05 list_users page=1 size=2 start=0 end=2 (输出可能变化/不固定：包含时间戳)
		data := ListUsersData{
			Items: users[start:end],
			Page:  page,
			Size:  size,
			Total: total,
		}
		writeJSON(w, http.StatusOK, APIResponse{Code: 0, Message: "OK", Data: data}) // Output: {"code":0,"message":"OK","data":{"items":[...],"page":1,"size":2,"total":4}}\n
	})

	addr := ":8080"
	log.Printf("listening on %s", addr) // Output: 2006/01/02 15:04:05 listening on :8080 (输出可能变化/不固定：包含时间戳)
	err := http.ListenAndServe(addr, mux)
	log.Printf("server stopped: %v", err) // Output: 2006/01/02 15:04:05 server stopped: listen tcp :8080: bind: address already in use (输出可能变化/不固定：取决于环境与错误)
}
```

### G. 怎么运行（命令 + 预期现象）
- 启动：`cd go-learning && go run ./cmd/day04_01b_json_errors`
- 预期现象：`curl` `/api/v1/health` 返回统一 JSON；`/api/v1/users` 支持分页，非法参数返回 400 + 业务码。

### H. 练习题（1–3 题，覆盖边界条件）

#### 练习 1：抽一个 `onlyGET` 复用（避免每个 handler 都手写 Method 判断）
- 验收标准：
  - `/api/v1/ping` 只支持 GET
  - `curl -s http://localhost:18080/api/v1/ping` 返回 `{"code":0,"message":"OK","data":{"pong":"pong"}}`
  - `curl -s -X POST http://localhost:18080/api/v1/ping` 返回 `{"code":10001,"message":"METHOD_NOT_ALLOWED"}`

#### 练习 2：把 404 也统一成 JSON（路径不存在也要一致结构）
- 验收标准：
  - `curl -s http://localhost:18080/api/v1/not-exist` 返回 `{"code":10004,"message":"NOT_FOUND"}`

### I. 参考答案（每题后紧跟可运行答案）

#### 练习 1 参考答案（可运行）
- 运行：`cd /Users/zhang/Desktop/go-study/codex/go-learning && PORT=18080 go run ./cmd/day04_01b_json_errors_ex1`
- 验证：
  - `curl -s http://localhost:18080/api/v1/ping` 典型响应：`{"code":0,"message":"OK","data":{"pong":"pong"}}`
  - `curl -s -X POST http://localhost:18080/api/v1/ping` 典型响应：`{"code":10001,"message":"METHOD_NOT_ALLOWED"}`

代码（全文）：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day04_01b_json_errors_ex1/main.go`

#### 练习 2 参考答案（可运行）
- 运行：`cd /Users/zhang/Desktop/go-study/codex/go-learning && PORT=18080 go run ./cmd/day04_01b_json_errors_ex2`
- 验证：
  - `curl -s http://localhost:18080/api/v1/health` 典型响应：`{"code":0,"message":"OK","data":{"ok":true}}`
  - `curl -s http://localhost:18080/api/v1/not-exist` 典型响应：`{"code":10004,"message":"NOT_FOUND"}`

代码（全文）：`/Users/zhang/Desktop/go-study/codex/go-learning/cmd/day04_01b_json_errors_ex2/main.go`

---

## References
- 官方：`net/http` 包文档 https://pkg.go.dev/net/http （Handler/ServeMux/Server 等概念）
- 官方：Go Blog/Effective Go https://go.dev/doc/effective_go （工程习惯与风格）
- 内部规范参考：`api-design-principles`（资源命名、HTTP 方法语义、错误响应一致性、分页等）
