package main

import "fmt"

func main() {
	fmt.Println("== Day02.4: defer (LIFO / args eval / named return) ==") // Output: == Day02.4: defer (LIFO / args eval / named return) ==

	demoLIFO()
	demoArgsEvaluatedAtDefer()

	fmt.Printf("namedReturn()=%d\n", namedReturn()) // Output: namedReturn()=2
}

func demoLIFO() {
	fmt.Println("demoLIFO start") // Output: demoLIFO start
	defer fmt.Println("defer 1")  // Output: defer 1
	defer fmt.Println("defer 2")  // Output: defer 2
	fmt.Println("demoLIFO end")   // Output: demoLIFO end
}

func demoArgsEvaluatedAtDefer() {
	fmt.Println("demoArgs start") // Output: demoArgs start
	x := 1
	defer fmt.Printf("defer x=%d\n", x) // Output: defer x=1
	x = 2
	fmt.Printf("now x=%d\n", x) // Output: now x=2
	fmt.Println("demoArgs end") // Output: demoArgs end
}

func namedReturn() (result int) {
	defer func() {
		result++
	}()
	return 1
}
