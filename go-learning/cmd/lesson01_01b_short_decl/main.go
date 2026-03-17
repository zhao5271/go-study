package main

import "fmt"

func main() {
	fmt.Println("== 01.1B := short declaration ==") // Output: == 01.1B := short declaration ==

	x := 42
	fmt.Printf("x=%d type=%T\n", x, x) // Output: x=42 type=int
}
