package main

import "fmt"

func main() {
	fmt.Println("== Lesson 01.1: var / := / zero values / type conversion ==") // Output: == Lesson 01.1: var / := / zero values / type conversion ==

	// 1) var + zero values
	// Why: Go doesn't have `undefined`. Declarations always produce a usable value (zero value).
	var count int
	var name string
	var ok bool
	fmt.Printf("count=%d name=%q ok=%v\n", count, name, ok) // Output: count=0 name="" ok=false

	// 2) short declaration := (function scope only)
	// Why: fast local development + type inference, but still static typing.
	x := 42
	fmt.Printf("x=%d type=%T\n", x, x) // Output: x=42 type=int

	// 3) explicit type when width matters
	var big int64 = 1
	fmt.Printf("big=%d type=%T\n", big, big) // Output: big=1 type=int64

	// 4) no implicit numeric conversions (TS/JS differs)
	// Go requires an explicit conversion.
	sum := int64(x) + big
	fmt.Printf("int64(x)+big=%d type=%T\n", sum, sum) // Output: int64(x)+big=43 type=int64

	// 5) constants are untyped until they need a type
	const pi = 3.14159
	fmt.Printf("pi=%.2f type=%T\n", pi, pi) // Output: pi=3.14 type=float64
}
