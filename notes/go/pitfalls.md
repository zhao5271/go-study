# Go 坑点库（TS/Node → Go 常见踩坑）

> 规则：每条尽量 1–3 行，附“怎么避免”。

## 大括号换行
- `{` 不能另起一行（自动插分号导致编译错）。
- 避免：按 gofmt 风格写，交给格式化器。

## 未使用即报错
- 未使用的 import/变量会编译失败。
- 避免：及时删除无用依赖；先写后用（或先用 `_ = x` 临时消除，但别滥用）。

## shadowing（:=）
- `:=` 在块内可能创建新变量，外层不更新。
- 避免：需要更新外层变量时用 `=`；在 if init 中刻意缩小作用域。

## iota 以“行”为单位递增
- 即使某行写了显式值，`iota` 也会继续按行 +1；每个 `const (...)` 块都会从 0 重置。
- 避免：把同一组枚举放在一个 const 块里；需要跳过某个值就用 `_ = iota` 明确占位。

## `_` 不要吞掉重要值
- `_` 会永久丢弃那个值；把 `error` 用 `_` 忽略会让线上问题更难定位。
- 避免：只对“确实不关心”的返回值用 `_`；对 error 要么处理，要么明确记录原因。

## 编译期接口校验不是“创建对象”
- `var _ IFace = (*T)(nil)` 只是编译期检查 `*T` 是否实现接口；不会在运行期构造对象。
- 避免：别把它当初始化/注册逻辑；需要注册请显式调用函数。

## `string(数字)` 不是数字字符串
- `string(65)` 得到 `"A"`（按字符码点转换），不是 `"65"`。
- 避免：数字转字符串用 `strconv.Itoa`；字符串转数字用 `strconv.Atoi`。

## 整数除法会截断
- `5/2 == 2`（不是 2.5）；如果你期望小数，必须先转成 `float64` 再除。
- 避免：`float64(a)/float64(b)`；或者明确说明你要的是“向下取整”的整数结果。

## float → int 会截断（不是四舍五入）
- 例如 `f := 1.9; int(f) == 1`，`f := -1.9; int(f) == -1`（朝 0 截断）。
- 避免：需要四舍五入就先 `math.Round`；需要保留小数就别转 int（用 `decimal` 类库/分为单位的整数）。

## `len(s)` 不是“字符数”
- Go 的 `string` 是 UTF-8 字节序列：`len(s)` 是字节数；`s[i]` 取到的是 `byte`，不是“第 i 个字符”。
- 避免：按字符处理用 `for range` 或 `utf8.RuneCountInString`。

## `int` 位宽相关（跨平台/交叉编译）
- 32 位平台 `int` 只有 32 位；`Atoi` 也只保证解析到 `int` 范围。
- 避免：需要稳定 64 位时用 `ParseInt(..., 64)`；`int64→int` 转换前做溢出检查。

## 位运算省括号导致误读/误写
- 权限判断要写 `(mask & PermRead) != 0`，不要写成 `mask & PermRead != 0`（可读性差且容易误解）。
- 避免：位运算一律加括号；mask/flag 用自定义类型提升可读性（`type Perm uint64`）。

## `%w` 滥用
- wrap 会暴露底层错误语义，可能破坏抽象边界。
- 避免：只有需要上层 `errors.Is/As` 判断时才 `%w`。

## slice 共享底层数组
- 子切片改元素会影响原切片；append 扩容后共享关系又变化。
- 避免：需要隔离就 `copy`；或三下标切片限制 cap。

## nil map 写入 panic
- `var m map[K]V` 没 make 时写入会 panic。
- 避免：用 `make(map[K]V)` 初始化。

## map 遍历顺序不稳定
- 输出/逻辑不要依赖顺序。
- 避免：取 keys 排序再遍历。

## defer 参数求值时机
- defer 的参数在 defer 语句那一刻就求值，不是执行时。
- 避免：需要“执行时值”就用闭包捕获（注意捕获变量的时机）。

## recover 作用域
- recover 只能在 defer 中、且只对当前 goroutine。
- 避免：不要把 recover 当通用异常处理；仅用于边界兜底。

## HTTP 响应：写了 body 就别再改 header/status
- `WriteHeader` 只能写一次；一旦写了 body，header/status 可能已经发出。
- 避免：先校验参数与权限；确定 status 后一次性写出 JSON。

## 405/404 结构不一致
- 成功是 JSON，失败却是纯文本/HTML，前端会被迫写分支。
- 避免：统一 `writeError`，尽早建立错误码与响应结构规范。

## MySQL init.sql 不会每次都执行
- Docker MySQL 的 `/docker-entrypoint-initdb.d` 只在“首次创建数据卷”时执行。
- 避免：改了 init.sql 后要 `docker compose down -v` 清卷再启动（注意会丢数据）。

## DSN 忘记 parseTime
- 扫描 `DATETIME/TIMESTAMP` 到 `time.Time` 可能失败或行为异常。
- 避免：MySQL DSN 加 `parseTime=true`。

## 值接收者“改了不生效”
- `struct` 是按值拷贝；值接收者方法修改的是副本。
- 避免：需要修改接收者时用指针接收者 `func (t *T) ...`；或返回新值并显式赋回。

## range 遍历 `[]struct` 修改不生效
- `for _, v := range s` 的 `v` 是副本，改字段不会回写到切片元素。
- 避免：用索引遍历 `for i := range s { s[i].Field = ... }` 或改成 `[]*T`。

## HTTP 工具函数复制粘贴导致不一致
- 同一个项目里出现多个 `writeJSON/writeError/parsePageSize` 版本，很容易某个接口漏校验/响应结构不同。
- 避免：抽到 `internal/httpkit`，所有 handler 只调用这一份实现。

## nil interface 坑（看起来不是 nil）
- `var repo UserRepo = (*MemoryUserRepo)(nil)`：`repo != nil` 但内部指针是 nil。
- 避免：接口里存指针时要小心判空；必要时显式检查底层指针或避免返回“带类型的 nil”。

## 方法接收者导致“没实现接口”
- `func (t *T) M()` 只在 `*T` 上；如果你用 `T{}` 去赋值接口，可能不满足方法集合。
- 避免：统一用指针接收者并传 `&T{}`；或确保值接收者满足你的接口设计。
