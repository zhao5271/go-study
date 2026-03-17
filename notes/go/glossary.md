# Go 术语表（外部记忆）

> 规则：每个术语 3–8 行；只写“能复用的关键点”，不要写长教程。

## Zero Value（零值）
- Go 没有 `undefined`；声明即初始化为类型的默认值（0/""/false/nil）。
- 好处：减少“未初始化就用”的运行期问题；代价：需要区分“零值就是业务有效值”与“缺省”。

## Exported / Unexported（导出/未导出）
- 标识符首字母大写：包外可见；小写：包内可见。
- 用命名约定替代 `public/private/export` 关键字。

## (value, error)
- 常见函数返回形式：成功返回 value + nil；失败返回零值 + error。
- 让错误处理显式化（对照 TS/Node try/catch）。

## Sentinel Error（哨兵错误）
- 用 `var ErrXxx = errors.New("...")` 表示可匹配语义。
- 匹配用 `errors.Is(err, ErrXxx)`（不要用 `==`，尤其当错误会 wrap）。

## Typed Error（类型错误）
- 自定义 error 类型携带结构化字段（如 Resource/ID）。
- 用 `errors.As(err, &target)` 提取。

## Wrap / %w
- `fmt.Errorf("...: %w", err)`：在加上下文的同时保留底层语义（可 `errors.Is/As`）。
- 取舍：wrap 会“暴露语义”给调用方，等于 API 承诺的一部分；不想暴露用 `%v`。

## Shadowing（变量遮蔽）
- `:=` 在块内可能创建同名新变量，外层变量不会被更新。
- 常见坑：`if _, err := f(); err != nil {}` 里的 `err` 不会影响外层 `err`。

## Defer（三条规则）
- 参数在 defer 语句执行时求值；defer 执行顺序 LIFO；可修改具名返回值。

## Panic / Recover
- 业务失败返回 error；panic 仅用于不可恢复的程序错误。
- recover 只能在 defer 中生效，且只对当前 goroutine。

## Slice backing array（slice 底层数组）
- slice 可能共享底层数组；子切片改元素会影响原切片。
- 需要隔离：`copy` 到新 slice 或用三下标切片限制 cap。

## nil vs empty（slice/map）
- nil slice JSON 通常是 `null`；empty slice JSON 是 `[]`（取决于编码场景与类型）。
- nil map 不能写入（会 panic），必须 `make(map[K]V)`。

## net/http Handler / ServeMux
- `http.Handler`：处理一次 HTTP 请求的抽象（有 `ServeHTTP(w, r)` 方法）。
- `http.HandlerFunc`：函数适配器，让普通函数也能当 Handler。
- `http.ServeMux`：最基础的路由分发器（按路径匹配把请求交给不同 Handler）。

## HTTP Status Code vs 业务错误码
- HTTP status：协议层语义（400/401/403/404/405/500…），便于网关/监控/客户端通用处理。
- 业务 code：业务层原因（如参数非法、鉴权失败、资源不存在），便于前端提示、埋点、告警聚合。
