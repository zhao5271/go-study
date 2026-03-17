# Go 工程模式库（外部记忆）

> 目的：把“可复用的工程套路”写下来，避免每次都从头推导。

## 错误处理（服务/仓库层）
- 函数签名优先：`func (s *Service) Do(ctx context.Context, ...) (T, error)`
- 模式：`v, err := ...; if err != nil { return zero, fmt.Errorf("context: %w", err) }`
- 是否 `%w`：只对需要上层判断语义的错误 wrap（否则用 `%v`）。

## HTTP Handler（提前约定）
- handler 只负责：解析输入 → 调 service → 把 error 映射成 HTTP 响应。
- 内部逻辑不要直接写响应；统一返回 error，在入口层转换。

## Table-driven tests 模板
```go
tests := []struct {
	name string
	in   int
	want int
	wantErrIs error
}{
	{"case1", 1, 2, nil},
}

for _, tt := range tests {
	t.Run(tt.name, func(t *testing.T) {
		got, err := Fn(tt.in)
		// 断言 got、errors.Is(err, tt.wantErrIs) ...
	})
}
```

## 项目结构（学习用最小版）
- `go-learning/cmd/dayNN_*`：每个知识点一个可运行入口
- `go-learning/internal/...`：可复用逻辑 + 可测试
- 笔记：`notes/go/dayNN-*.md` + 索引文件（progress/glossary/patterns/pitfalls）

