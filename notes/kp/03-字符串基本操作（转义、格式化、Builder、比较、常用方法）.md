---
type: kp
domain: go
topic: string-basics
topic_zh: 字符串基本操作
stage: foundation
status: evergreen
review_cycle: weekly
source:
  - official
tags:
  - go
  - kp
  - foundation
created: 2026-03-20
updated: 2026-03-20
---

# 03 字符串基本操作（转义、格式化、`strings.Builder`、比较、常用方法）

> 更新于：2026-03-20  
> 目标：把字符串在 Go 里的“表示、格式化、拼接、比较、常用处理”整理成后端项目里能直接套用的知识卡片。

## TL;DR（可放入 progress/context-pack）
- Go 的 `string` 是不可变的 UTF-8 字节序列；`len(s)` 是字节数，不是字符数。
- 字符串字面量两种：`"..."` 支持转义，`` `...` `` 原样保留且可多行。
- `fmt.Sprintf` 适合拼日志/错误信息；`%q` 很适合排查字符串里的换行和转义。
- 循环拼接优先用 `strings.Builder`；已知大概长度时先 `Grow`。
- 相等比较直接用 `==`；大小写不敏感时优先 `strings.EqualFold`。
- `Trim` 是按字符集合修剪；删固定前缀/后缀用 `TrimPrefix` / `TrimSuffix`。

## 关键词
- string literal
- raw / interpreted string
- UTF-8 / byte / rune
- Sprintf
- strings.Builder
- EqualFold / TrimPrefix / Split / Join

## 知识点 1：字符串字面量、转义、字节与 rune

### 一句话定义
Go 的字符串是不可变的字节序列；你看到的“字符”在底层往往是多个 UTF-8 字节。

### 为什么重要
后台管理 API 经常处理用户名、搜索词、日志文本、JSON 片段、SQL 片段；如果你把“字节数”和“字符数”混为一谈，截断、校验、展示都会出问题。

### 重难点拆解
- `"...“` 是解释型字符串，支持 `\n`、`\t`、`\"` 等转义。
- `` `...` `` 是原生字符串，反斜杠没有特殊含义，适合多行文本、正则、SQL。
- `len(s)` 看的是字节数；按字符遍历用 `for range` 或 `utf8.RuneCountInString`。
- `s[i]` 取到的是 `byte`，不是“第 i 个字符”。

### 业务场景落地
- 搜索框长度校验：你要先明确你限制的是字节还是字符
- 日志调试：`%q` 能把隐藏换行、空格、转义字符显示出来
- 写 SQL/JSON 模板时，原生字符串通常更可读

### 代码示例
```go
package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	interp := "a\nb\tc\\\""
	fmt.Printf("interp=%q\n", interp) // Output: interp="a\nb\tc\\\""

	raw := `a\nb\tc\\\"`
	fmt.Printf("raw=%q\n", raw) // Output: raw="a\\nb\\tc\\\\\\\""

	s := "中文"
	fmt.Printf("len_bytes=%d\n", len(s))                 // Output: len_bytes=6
	fmt.Printf("rune_count=%d\n", utf8.RuneCountInString(s)) // Output: rune_count=2
}
```

### 怎么运行
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/kp/string-basics
# 预期会看到：
# == KP: string basics ==
# [1] literals + escapes
# interp="a\nb\tc\\\""
# raw="a\\nb\\tc\\\\\\\""
```

### 练习题 1
为什么 `len("中A") == 4`，但“字符数”是 2？

**验收标准**
- 能说明 `len` 统计的是字节数
- 能说明 UTF-8 下中文通常占多个字节

**参考答案**
- 因为 Go 的 `string` 底层是 UTF-8 字节序列；`"中"` 占 3 个字节，`"A"` 占 1 个字节，所以总字节数是 4。
- 如果要看“字符数”，应该用 `for range` 或 `utf8.RuneCountInString("中A")`，结果是 2。

## 知识点 2：格式化输出与 `strings.Builder`

### 一句话定义
`fmt` 负责格式化，`strings.Builder` 负责高效拼接字符串。

### 为什么重要
日志、错误信息、审计记录、动态 SQL 片段、导出文本都会频繁拼接字符串；简单场景可以直接 `Sprintf`，循环或多段拼接更适合 `Builder`。

### 重难点拆解
- `fmt.Sprintf` 返回字符串，适合“先拼好，再交给日志/响应层”。
- `%q` 会带引号并转义，适合调试；`%x` 常用于十六进制输出。
- `strings.Builder` 的零值可直接使用；非零 Builder 不要复制。
- 已知大概长度时先 `Grow`，可以减少扩容和拷贝。

### 业务场景落地
- 审计日志：`user=alice action=create_user status=200`
- 动态 WHERE：只拼 SQL 结构，值仍然走参数化
- 导出文本/CSV：多段拼接时 Builder 更稳

### 代码示例
```go
package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	msg := fmt.Sprintf("audit user=%s status=%d", "alice", 200)
	fmt.Printf("msg=%q\n", msg) // Output: msg="audit user=alice status=200"

	var b strings.Builder
	b.Grow(64)
	b.WriteString("audit user=")
	b.WriteString("alice")
	b.WriteString(" status=")
	b.WriteString(strconv.Itoa(200))
	fmt.Printf("%s\n", b.String()) // Output: audit user=alice status=200
}
```

### 怎么运行
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/kp/string-basics
# 预期会看到：
# [2] fmt formatting
# msg="audit user=alice status=200"
# [3] strings.Builder
# audit user=alice action=CREATE_USER status=200
```

### 练习题 2
把下面的循环 `+=` 改成 `strings.Builder`：

```go
line := ""
for _, part := range []string{"user=alice", "action=create"} {
	line += part
}
```

**验收标准**
- 使用 `strings.Builder`
- 最终结果是 `user=aliceaction=create` // Output: user=aliceaction=create

**参考答案**
```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	var b strings.Builder
	for _, part := range []string{"user=alice", "action=create"} {
		b.WriteString(part)
	}
	fmt.Printf("%s\n", b.String()) // Output: user=aliceaction=create
}
```

## 知识点 3：字符串比较与常用方法

### 一句话定义
比较和处理字符串时，要先想清楚：你要的是“完全相等”、还是“大小写无关”、还是“删前缀/按分隔符拆分”。

### 为什么重要
角色名、状态值、搜索条件、CSV 字段、路径前缀、批量导出字段拼接，这些都离不开字符串比较与处理。

### 重难点拆解
- 完全相等直接用 `==`；大小写不敏感用 `strings.EqualFold`。
- `Trim` 是按字符集合修剪，不是删固定子串；固定前缀/后缀请用 `TrimPrefix/TrimSuffix`。
- `Split("", ",")` 的结果不是空切片，而是 `[]string{""}`。
- 对用户输入常见的第一步是 `TrimSpace`，再做比较或分割。

### 业务场景落地
- 角色名比较：`strings.EqualFold(role, "admin")`
- 标签/导出字段解析：`Split + TrimSpace + Join`
- API 路径清洗：`TrimPrefix`

### 代码示例
```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Printf("eq=%v\n", "admin" == "admin")                    // Output: eq=true
	fmt.Printf("fold=%v\n", strings.EqualFold("Admin", "admin")) // Output: fold=true

	q := "  admin , editor  "
	q = strings.TrimSpace(q)
	fmt.Printf("q=%q\n", q) // Output: q="admin , editor"

	parts := strings.Split(q, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	fmt.Printf("parts=%v\n", parts) // Output: parts=[admin editor]
	fmt.Printf("prefix=%q\n", strings.TrimPrefix("/api/v1/users", "/api")) // Output: prefix="/v1/users"
}
```

### 怎么运行
```bash
cd /Users/zhang/Desktop/go-study/codex/go-learning
go run ./cmd/kp/string-basics
# 预期会看到：
# [4] compare
# eq=true
# fold=true
# [5] common strings methods
# trimmed=[admin editor]
```

### 练习题 3
把输入 `" Admin "` 处理成：
- 去掉首尾空格
- 与 `"admin"` 做大小写无关比较

**验收标准**
- 使用 `TrimSpace`
- 使用 `EqualFold`
- 输出 `ok=true` // Output: ok=true

**参考答案**
```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	s := strings.TrimSpace(" Admin ")
	ok := strings.EqualFold(s, "admin")
	fmt.Printf("ok=%t\n", ok) // Output: ok=true
}
```

## 关联复习
- `02-基础类型（转换、格式化、表达式）` 里已经讲过 `fmt` / `strconv` 的职责划分。
- `glossary.md` 里保留了 `rune / byte`、raw/interpreted string 的最小定义，方便快速查。

## References
- [Go Spec](https://go.dev/ref/spec) - 字符串类型、raw/interpreted string literal、转义语义的官方定义（官方）
- [strings package](https://pkg.go.dev/strings) - `Builder`、`EqualFold`、`TrimPrefix`、`Split`、`Join` 等标准库 API（官方）
- [unicode/utf8 package](https://pkg.go.dev/unicode/utf8) - `RuneCountInString` 等 UTF-8 工具函数（官方）
- NotebookLM 查询失败（原因：浏览器 profile 被占用，`ProcessSingletonLock` 冲突），本节结论以官方资料补齐。
