package main

import "fmt"

func main() {
	fmt.Println("== 01.1C explicit conversion (int vs int64) ==") // Output: == 01.1C explicit conversion (int vs int64) ==

	x := 42
	var big int64 = 1
	sum := int64(x) + big
	fmt.Printf("int64(x)+big=%d type=%T\n", sum, sum) // Output: int64(x)+big=43 type=int64
}
