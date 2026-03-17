package main

import "fmt"

func main() {
	fmt.Println("== Day02.1: functions + multiple returns + named returns ==") // Output: == Day02.1: functions + multiple returns + named returns ==

	x, y := split(10)
	fmt.Printf("split(10) => x=%d y=%d\n", x, y) // Output: split(10) => x=4 y=6

	q, r := divmod(17, 5)
	fmt.Printf("divmod(17,5) => q=%d r=%d\n", q, r) // Output: divmod(17,5) => q=3 r=2

	fmt.Printf("namedZero()=%d\n", namedZero()) // Output: namedZero()=0
}

// split demonstrates named return values. They are real variables (zero-valued).
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}

func divmod(a, b int) (int, int) {
	return a / b, a % b
}

func namedZero() (x int) {
	return
}
