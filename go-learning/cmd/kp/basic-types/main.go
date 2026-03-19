package main

import (
	"fmt"
	"math"
	"strconv"
	"unicode/utf8"
)

func main() {
	fmt.Println("== KP: basic types ==") // Output: == KP: basic types ==

	demoConversions()
	demoFormatting()
	demoExpressions()
}

func demoConversions() {
	fmt.Println("[1] conversions") // Output: [1] conversions

	x := 42
	var big int64 = 1
	sum := int64(x) + big
	fmt.Printf("sum=%d type=%T\n", sum, sum) // Output: sum=43 type=int64

	// int division truncates; convert to float to keep decimals.
	ratio := float64(3) / 2
	fmt.Printf("float64(3)/2=%.2f\n", ratio) // Output: float64(3)/2=1.50

	fmt.Printf("IntSize=%d\n", strconv.IntSize) // Output: IntSize=64（输出可能变化/不固定：取决于运行平台是 32/64 位）

	s := "中A"
	fmt.Printf("len(bytes)=%d\n", len(s))                    // Output: len(bytes)=4
	fmt.Printf("runes=%d\n", utf8.RuneCountInString(s))       // Output: runes=2
	fmt.Printf("firstByte=%d\n", s[0])                        // Output: firstByte=228（说明：这是 UTF-8 首字节，不等于字符）
	fmt.Printf("firstRune=%q\n", []rune(s)[0])                // Output: firstRune='中'
	f1 := 1.9
	fmt.Printf("int(1.9)=%d\n", int(f1))                 // Output: int(1.9)=1
	fmt.Printf("round(1.9)=%d\n", int(math.Round(f1)))   // Output: round(1.9)=2
	f2 := -1.9
	fmt.Printf("int(-1.9)=%d\n", int(f2))                // Output: int(-1.9)=-1
	fmt.Printf("round(-1.9)=%d\n", int(math.Round(f2)))  // Output: round(-1.9)=-2
}

func demoFormatting() {
	fmt.Println("[2] formatting + strconv") // Output: [2] formatting + strconv

	n := 65
	fmt.Printf("n=%d\n", n) // Output: n=65

	// string(n) treats n as a rune/byte -> "A", not "65".
	fmt.Printf("string(n)=%q\n", string(n)) // Output: string(n)="A"
	fmt.Printf("Itoa(n)=%q\n", strconv.Itoa(n)) // Output: Itoa(n)="65"

	v, err := strconv.Atoi("65")
	if err != nil {
		fmt.Printf("Atoi error=%v\n", err) // Output: （输出可能变化/不固定：错误信息随 Go 版本/实现变化）
		return
	}
	fmt.Printf("Atoi(\"65\")=%d\n", v) // Output: Atoi("65")=65

	v64, err := strconv.ParseInt("9223372036854775807", 10, 64)
	if err != nil {
		fmt.Printf("ParseInt error=%v\n", err) // Output: （输出可能变化/不固定：错误信息随输入/实现变化）
		return
	}
	fmt.Printf("v64=%d type=%T\n", v64, v64) // Output: v64=9223372036854775807 type=int64

	raw := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	fmt.Printf("raw=%x\n", raw) // Output: raw=deadbeef

	str := "hello"
	fmt.Printf("str=%q type=%T\n", str, str) // Output: str="hello" type=string
}

func demoExpressions() {
	fmt.Println("[3] operators/expressions") // Output: [3] operators/expressions

	fmt.Printf("5/2=%d\n", 5/2) // Output: 5/2=2
	fmt.Printf("5/2=%.1f\n", float64(5)/2) // Output: 5/2=2.5

	var p *int
	if p != nil && *p > 0 {
		fmt.Println("p>0") // Output: （不会打印，因为 p==nil，短路）
	} else {
		fmt.Println("short-circuit ok (p is nil)") // Output: short-circuit ok (p is nil)
	}

	offset, err := calcOffset(2, 10)
	if err != nil {
		fmt.Printf("offset err=%v\n", err) // Output: （输出可能变化/不固定：错误信息随输入/实现变化）
		return
	}
	fmt.Printf("offset(page=2,size=10)=%d\n", offset) // Output: offset(page=2,size=10)=10

	mask := PermRead | PermExport
	fmt.Printf("mask=%b\n", mask) // Output: mask=101
	fmt.Printf("hasRead=%t\n", hasPerm(mask, PermRead))   // Output: hasRead=true
	fmt.Printf("hasWrite=%t\n", hasPerm(mask, PermWrite)) // Output: hasWrite=false
}

func calcOffset(page, size int) (int, error) {
	if page < 1 {
		return 0, fmt.Errorf("page must be >= 1")
	}
	if size < 1 || size > 100 {
		return 0, fmt.Errorf("size must be 1..100")
	}
	return (page - 1) * size, nil
}

type Perm uint64

const (
	PermRead Perm = 1 << iota
	PermWrite
	PermExport
)

func hasPerm(mask, p Perm) bool { return (mask&p) != 0 }
