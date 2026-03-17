package main

import "fmt"

func main() {
	fmt.Println("== 01.1A var + zero values ==") // Output: == 01.1A var + zero values ==

	var count int
	var name string
	var ok bool
	fmt.Printf("count=%d name=%q ok=%v\n", count, name, ok) // Output: count=0 name="" ok=false
}
