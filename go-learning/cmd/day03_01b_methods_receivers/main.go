package main

import "fmt"

type Counter struct {
	N int
}

// 值接收者：操作的是“副本”（copy）
func (c Counter) IncByValue() {
	c.N++
}

// 指针接收者：操作的是“同一个对象”
func (c *Counter) IncByPtr() {
	c.N++
}

func main() {
	fmt.Println("== Day03.1b: methods + value/pointer receiver ==") // Output: == Day03.1b: methods + value/pointer receiver ==

	c := Counter{N: 10}
	c.IncByValue()
	fmt.Printf("after IncByValue: c.N=%d\n", c.N) // Output: after IncByValue: c.N=10

	c.IncByPtr()
	fmt.Printf("after IncByPtr:   c.N=%d\n", c.N) // Output: after IncByPtr:   c.N=11

	// Go 会在需要时自动取地址/解引用，让方法调用更顺滑
	cp := &Counter{N: 20}
	cp.IncByPtr()
	fmt.Printf("cp.N=%d\n", cp.N) // Output: cp.N=21
}

