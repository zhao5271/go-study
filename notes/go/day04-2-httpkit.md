# Day 04.2：抽可复用 HTTP 工具包（`writeJSON/writeError/parsePageSize`）——避免复制粘贴

> 贯穿项目：后台管理 API（RBAC + 登录鉴权 + 列表分页检索 + CRUD + 审计日志 + Docker 部署）  
> 本节目标：把统一 JSON 响应与分页 query 解析抽成可复用包，让后续 Gin/更多接口不会因为复制粘贴而“响应结构不一致/校验不一致”。

统一运行目录：
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
```

---

## 知识点 1：抽 `writeJSON/writeError`（统一响应结构的“唯一入口”）

### B. 一句话定义
把“写 JSON + 写统一错误响应”封装成函数，所有 handler 都只调用它，避免每个文件各写各的。

### C. 为什么重要（不做会怎样）
后台管理 API 一旦接口多了，你会遇到：
- 有的接口返回 `{"code":...}`，有的直接 `http.Error` 文本
- 有的接口 Content-Type 不对
最终前端/调用方被迫写大量兼容分支，排障与维护成本飙升。

### D. 重难点拆解（2–4 条）
1) **统一输出口**：越早抽出来越好，后面改字段/加 request id 更容易。  
2) **不要在业务逻辑里直接写响应**：业务只返回结果/错误，入口层统一写响应。  
3) **输出标注规则**：工具函数里的输出取决于入参，不可能写死具体 JSON，需要标注“输出可能变化/不固定”。

### E. 业务场景落地（后台管理 API）
- 所有接口返回结构统一：`{code,message,data}`  
- 所有错误统一走 `WriteError`，405/404/400/500 结构一致。

### F. 代码示例（最小可运行）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/internal/httpkit/response.go`
```go
package httpkit

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v) // Output: 输出可能变化/不固定（取决于 v 的内容）
}

func WriteError(w http.ResponseWriter, status int, code int, message string) {
	WriteJSON(w, status, APIResponse{Code: code, Message: message})
}
```

### G. 怎么运行（命令 + 预期现象）
这个文件是库文件，不能直接 `go run`。你可以通过现有示例验证它生效：
```bash
PORT=18080 go run ./cmd/day04_01b_json_errors
# Output: 2006/01/02 15:04:05 listening on :18080 (输出可能变化/不固定：包含时间戳)
```
然后：
```bash
curl -s http://localhost:18080/api/v1/health
# Output: {"code":0,"message":"OK","data":{"ok":true}}
```

### H. 练习题（1–3 题，覆盖边界条件）
练习 1：把你自己的新 handler（比如 `/api/v1/ping`）也改成只用 `httpkit.WriteJSON/WriteError`  
- 验收标准：成功与失败都返回统一 JSON（用 curl 验证）

### I. 参考答案（紧跟练习题）
参考答案 1：见 `day04_01b_json_errors_ex1`，它已经使用 `httpkit.WriteJSON/WriteError`：
```bash
PORT=18080 go run ./cmd/day04_01b_json_errors_ex1
# Output: 2006/01/02 15:04:05 listening on :18080 (输出可能变化/不固定：包含时间戳)
```

---

## 知识点 2：抽 `parsePageSize`（分页参数校验的“唯一真相”）

### B. 一句话定义
把 `page/size` 的默认值、范围校验与解析集中在一个函数里，所有列表接口都复用它。

### C. 为什么重要（不做会怎样）
列表接口最容易被滥用（`size=100000`），如果每个 handler 自己写校验，很容易某个接口漏了边界，拖垮服务。

### D. 重难点拆解（2–4 条）
1) **边界要统一**：例如 `page>=1`、`1<=size<=100`。  
2) **错误语义要稳定**：用一个 sentinel error（`ErrInvalidQuery`）表达“参数非法”，便于上层统一映射成 HTTP 400。  
3) **只做“解析+校验”**：不要在这里做业务逻辑。

### E. 业务场景落地（后台管理 API）
`GET /api/v1/users?page=1&size=20`、`GET /api/v1/audit-logs?page=1&size=20` 都复用同一套规则。

### F. 代码示例（最小可运行）
文件：`/Users/zhang/Desktop/go-study/codex/go-learning/internal/httpkit/query.go`
```go
package httpkit

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var ErrInvalidQuery = errors.New("invalid query")

func ParsePageSize(r *http.Request) (page int, size int, err error) {
	page = 1
	size = 20

	q := r.URL.Query()
	if raw := strings.TrimSpace(q.Get("page")); raw != "" {
		n, convErr := strconv.Atoi(raw)
		if convErr != nil || n < 1 {
			return 0, 0, ErrInvalidQuery
		}
		page = n
	}
	if raw := strings.TrimSpace(q.Get("size")); raw != "" {
		n, convErr := strconv.Atoi(raw)
		if convErr != nil || n < 1 || n > 100 {
			return 0, 0, ErrInvalidQuery
		}
		size = n
	}
	return page, size, nil
}
```

### G. 怎么运行（命令 + 预期现象）
用 `day04_01b_json_errors` 或 `day05_04_http_users_list` 验证解析与 400 返回：
```bash
PORT=18080 go run ./cmd/day04_01b_json_errors
# Output: 2006/01/02 15:04:05 listening on :18080 (输出可能变化/不固定：包含时间戳)
```
然后：
```bash
curl -s "http://localhost:18080/api/v1/users?page=0"
# Output: {"code":10002,"message":"INVALID_QUERY"}
```

### H. 练习题（1–3 题，覆盖边界条件）
练习 1：在一个新列表接口里也复用 `ParsePageSize`（比如 `/api/v1/audit-logs`，先返回空数组即可）  
- 验收标准：
  - `page=0` 返回 400 + `INVALID_QUERY`
  - `size=101` 返回 400 + `INVALID_QUERY`

### I. 参考答案
参考答案 1：照抄 `day04_01b_json_errors` 的 `/api/v1/users` 写法即可——解析失败就 `httpkit.WriteError(400,10002,"INVALID_QUERY")`。

---

## References
- 官方：Go `net/http` https://pkg.go.dev/net/http
- 官方：Go `encoding/json` https://pkg.go.dev/encoding/json
