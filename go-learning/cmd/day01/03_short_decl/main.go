package main

import "fmt"

func main() {
	fmt.Println("== Day01.3: := short declaration ==") // Output: == Day01.3: := short declaration ==

	x := 42
	fmt.Printf("x=%d type=%T\n", x, x) // Output: x=42 type=int

	// := is function-scope only; package-level must use var.
}
