package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("== Day01.1: package/import + exported names ==") // Output: == Day01.1: package/import + exported names ==

	s := "go fullstack"
	fmt.Printf("ToUpper(%q)=%q\n", s, strings.ToUpper(s)) // Output: ToUpper("go fullstack")="GO FULLSTACK"

	// In Go: identifiers starting with Uppercase are exported (public).
	// identifiers starting with lowercase are unexported (package-private).
}
