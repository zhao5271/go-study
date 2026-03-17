package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("== Day01.6: if/for/switch ==") // Output: == Day01.6: if/for/switch ==

	// if with init statement
	if n := len("go"); n > 1 {
		fmt.Printf("len=%d\n", n) // Output: len=2
	}

	// for (classic)
	fmt.Print("classic: ") // Output: classic:
	for i := 0; i < 3; i++ {
		fmt.Printf("%d ", i) // Output: 0 1 2
	}
	fmt.Println() // Output: (newline)

	// for (while-like)
	fmt.Print("while-like: ") // Output: while-like:
	x := 0
	for x < 3 {
		fmt.Printf("%d ", x) // Output: 0 1 2
		x++
	}
	fmt.Println() // Output: (newline)

	// switch (default break)
	var sb strings.Builder
	switch v := 2; v {
	case 1:
		sb.WriteString("one")
	case 2:
		sb.WriteString("two")
		fallthrough
	case 3:
		sb.WriteString("+three")
	default:
		sb.WriteString("other")
	}
	fmt.Printf("switch=%s\n", sb.String()) // Output: switch=two+three
}
