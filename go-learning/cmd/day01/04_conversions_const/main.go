package main

import "fmt"

func main() {
	fmt.Println("== Day01.4: explicit conversions + const ==") // Output: == Day01.4: explicit conversions + const ==

	x := 42
	var big int64 = 1
	sum := int64(x) + big
	fmt.Printf("int64(x)+big=%d type=%T\n", sum, sum) // Output: int64(x)+big=43 type=int64

	const pi = 3.14159
	fmt.Printf("pi=%.2f type=%T\n", pi, pi) // Output: pi=3.14 type=float64
}
