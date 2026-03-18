package main

import (
	"fmt"

	"example.com/go-learning/internal/day02/users"
)

func main() {
	fmt.Println("== Day02.3: if init + := shadowing ==") // Output: == Day02.3: if init + := shadowing ==

	user, err := users.FindUserSentinel(1)
	fmt.Printf("outer: user.ID=%d err=%v\n", user.ID, err) // Output: outer: user.ID=1 err=<nil>

	// Common pitfall: := inside if creates a NEW err that does not update outer err.
	if _, err := users.FindUserSentinel(2); err != nil {
		fmt.Printf("inside if: err=%v\n", err) // Output: inside if: err=find user 2: user not found
	}
	fmt.Printf("after if: outer err=%v\n", err) // Output: after if: outer err=<nil>

	// If you really want to update outer err, use assignment (=), not :=.
	_, err = users.FindUserSentinel(2)
	fmt.Printf("after assignment: err=%v\n", err) // Output: after assignment: err=find user 2: user not found
}
