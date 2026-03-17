package main

import "fmt"

func main() {
	fmt.Println("== Day01.2: var + zero values ==") // Output: == Day01.2: var + zero values ==

	var count int
	var name string
	var ok bool
	fmt.Printf("count=%d name=%q ok=%v\n", count, name, ok) // Output: count=0 name="" ok=false

	var p *int
	fmt.Printf("p==nil? %v\n", p == nil) // Output: p==nil? true
}
