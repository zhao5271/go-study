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

## DSN（MySQL 连接串）
- 连接数据库的字符串（用户名/密码/地址/库名/参数）。
- 本仓库示例：`app:app@tcp(127.0.0.1:3307)/go_admin?parseTime=true`（3307 是 compose 映射端口）。
- `parseTime=true`：让 `DATETIME/TIMESTAMP` 能正确扫描到 `time.Time`。

## Transaction（事务）
- 多步写入要么全成功，要么全失败；失败要回滚。
- 常见场景：创建用户 + 写审计日志，必须保证一致性。

## Receiver（方法接收者）
- 方法是“带接收者的函数”：`func (r T) M()` 或 `func (r *T) M()`。
- 值接收者拿到副本，通常不能修改原对象；需要修改时用指针接收者。
- 经验法则：业务服务/仓库类型一般用指针接收者。

## `omitempty`（JSON tag）
- `json:"field,omitempty"`：当字段是“空值”时序列化会省略该字段（具体空值语义随类型不同）。
- 常见用法：指针字段为 `nil` 时省略，用于表达 PATCH 的 optional 字段“没传”。

## Interface（接口）
- 接口是一组方法的契约；谁的方法集合匹配，谁就实现（隐式实现）。
- 工程上常用于分层边界：service 依赖 repo 接口，repo 可替换（内存/MySQL/mock）。

## Embedding（匿名嵌入/组合）
- 把一个类型匿名放进 struct：`type X struct { Y }`，Y 的方法会被提升（promoted）。
- 常用于装饰器：`type LoggingRepo struct{ Repo }`，在不改业务逻辑的情况下加 log/metrics。

## internal 包（internal/）
- `internal/` 下的包只能被“同一个 module 内部”导入，避免被外部项目误用。
- 工程上常用来放：可复用但不想对外承诺稳定 API 的实现细节（例如 `internal/httpkit`）。

## Scope（作用域 / block scope）
- 变量只在它声明所在的“作用域”内可见；最常见边界是大括号 `{}`（块级作用域）。
- 常见场景：`if/for/switch` 的 init 语句可缩小变量作用域，减少误用。

## Blank identifier（匿名变量 `_`）
- `_` 用来显式丢弃你不需要的值（例如忽略多返回值中的某一个、或 `range` 的 index/value）。
- 不要用 `_` 吞掉重要信息（尤其是 error）；它表示“我明确不关心”。

## iota（常量计数器）
- `iota` 只在 `const (...)` 块内生效：从 0 开始，每行 +1；每个 const 块会重置。
- 常用于：枚举常量（加 `+1` 跳过 0）、位标志（`1 << (10*iota)`）。

## Conversion（显式类型转换）
- Go 不做隐式数值类型转换；不同数值类型相加/比较需要显式转换（例如 `int64(x)`）。
- 转换可能导致溢出/截断（`float64→int` 截断小数；`int64→int` 可能溢出）。

## fmt verbs（格式化动词）
- 常用：`%d`（整数）、`%f`（浮点）、`%s`（字符串）、`%q`（带引号字符串）、`%v`（默认）、`%T`（类型）。
- 用 `%T` 快速确认值的真实类型，减少“类型以为对了”的调试时间。

## strconv（字符串转换）
- `strconv.Atoi("123")`：string→int；`strconv.Itoa(123)`：int→string。
- 别把 `string(123)` 当成 `"123"`：它会把整数当字符码点转换（通常得到单字符）。

## rune / byte（字符码点 / 字节）
- `rune` = `int32`：表示一个 Unicode code point（更像“字符”）。
- `byte` = `uint8`：表示一个字节（更像“原始数据”）。
- `string` 是 UTF-8 字节序列：`len(s)` 是字节数；`s[i]` 取到的是 `byte`；`for range s` 才是按 `rune` 迭代。

## unicode/utf8（UTF-8 工具包）
- `utf8.RuneCountInString(s)`：统计 rune 数（“字符数”）。
- `utf8.DecodeRuneInString(s)`：从字符串解码第一个 rune（做前缀解析时常用）。

## Bit flag（位标志）
- 用位运算把“多个布尔开关/权限”压到一个整数里：`mask := Read | Export`。
- 判断用 `(mask & Read) != 0`（括号不要省，避免误读/写错）。

## Untyped constant（无类型常量）
- `const x = 1` 在赋值前没有固定类型，会按上下文“落到”某个具体类型（只要能表示）。
- 一旦变成变量（例如 `var x = 1`），它就有了固定类型（通常是 `int`），跨类型运算必须显式转换。

## Comparable（可比较）
- 可用 `==/!=` 比较的类型称为 comparable（大部分基础类型都可以）。
- `slice/map/func` 不能互相 `==`（只能与 `nil` 比），工程上常用 `len(s)==0` 或 `maps.Equal` 等替代手段（后续再展开）。
