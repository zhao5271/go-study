package main

import (
	"fmt"

	"example.com/go-learning/internal/day02/users"
)

func main() {
	fmt.Println("== Day02.5: panic vs error + recover ==") // Output: == Day02.5: panic vs error + recover ==

	// Business logic failures should return error, not panic.
	_, err := users.FindUserSentinel(2)
	fmt.Printf("business err=%v\n", err) // Output: business err=find user 2: user not found

	// Panic is for programmer errors / truly unrecoverable situations.
	fmt.Println("panic demo start")                              // Output: panic demo start
	fmt.Printf("recovered=%v\n", safe(func() { panic("boom") })) // Output: recovered=boom
	fmt.Printf("recovered(no panic)=%v\n", safe(func() {}))      // Output: recovered(no panic)=<nil>
	fmt.Println("still running")                                 // Output: still running
}

func safe(fn func()) (recovered any) {
	defer func() {
		recovered = recover()
	}()
	fn()
	return nil
}
