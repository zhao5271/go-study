package main

import "fmt"

func main() {
	fmt.Println("== KP: vars/const/scope/_/iota ==") // Output: == KP: vars/const/scope/_/iota ==

	demoScope()
	fmt.Println("----") // Output: ----

	demoBlankIdentifier()
	fmt.Println("----") // Output: ----

	demoIota()
}

func demoScope() {
	fmt.Println("[1] scope + block scope + := shadowing") // Output: [1] scope + block scope + := shadowing

	x := 1
	if true {
		x := 2
		fmt.Printf("inner x=%d\n", x) // Output: inner x=2
	}
	fmt.Printf("outer x=%d\n", x) // Output: outer x=1

	// Tip: if you want to update the outer x, use `x = 2` (not `:=`) inside the block.
}

func demoBlankIdentifier() {
	fmt.Println("[2] blank identifier (_)") // Output: [2] blank identifier (_)

	q, _ := divMod(10, 3)
	fmt.Printf("10/3 quotient=%d\n", q) // Output: 10/3 quotient=3

	names := []string{"alice", "bob"}
	for _, name := range names {
		fmt.Printf("name=%s\n", name) // Output: name=alice (第二行 Output: name=bob)
	}
}

func divMod(a, b int) (q int, r int) {
	return a / b, a % b
}

type Role int

const (
	RoleAdmin Role = iota + 1
	RoleEditor
	RoleViewer
)

const (
	_ = iota // skip 0
	KB = 1 << (10 * iota)
	MB
	GB
)

func demoIota() {
	fmt.Println("[3] const + iota") // Output: [3] const + iota

	fmt.Printf("RoleAdmin=%d RoleEditor=%d RoleViewer=%d\n", RoleAdmin, RoleEditor, RoleViewer) // Output: RoleAdmin=1 RoleEditor=2 RoleViewer=3
	fmt.Printf("KB=%d MB=%d GB=%d\n", KB, MB, GB)                                              // Output: KB=1024 MB=1048576 GB=1073741824
}

