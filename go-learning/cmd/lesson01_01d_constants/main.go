package main

import "fmt"

func main() {
	fmt.Println("== 01.1D constants ==") // Output: == 01.1D constants ==

	const pi = 3.14159
	fmt.Printf("pi=%.2f type=%T\n", pi, pi) // Output: pi=3.14 type=float64
}
