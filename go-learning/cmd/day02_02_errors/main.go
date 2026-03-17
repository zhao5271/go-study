package main

import (
	"errors"
	"fmt"

	"example.com/go-learning/internal/day02/users"
)

func main() {
	fmt.Println("== Day02.2: (value, error) + wrap + errors.Is/As ==") // Output: == Day02.2: (value, error) + wrap + errors.Is/As ==

	_, err := users.FindUserSentinel(2)
	fmt.Printf("sentinel err=%v\n", err) // Output: sentinel err=find user 2: user not found

	fmt.Printf("errors.Is(err, ErrUserNotFound)=%v\n", errors.Is(err, users.ErrUserNotFound)) // Output: errors.Is(err, ErrUserNotFound)=true
	fmt.Printf("errors.Unwrap(err)=%v\n", errors.Unwrap(err))                                 // Output: errors.Unwrap(err)=user not found

	hidden := fmt.Errorf("hide: %v", users.ErrUserNotFound)
	fmt.Printf("hidden=%v\n", hidden)                               // Output: hidden=hide: user not found
	fmt.Printf("errors.Unwrap(hidden)=%v\n", errors.Unwrap(hidden)) // Output: errors.Unwrap(hidden)=<nil>

	_, err = users.FindUserTyped(2)
	fmt.Printf("typed err=%v\n", err) // Output: typed err=find user 2: user 2 not found

	var nf *users.NotFoundError
	fmt.Printf("errors.As(err, *NotFoundError)=%v\n", errors.As(err, &nf)) // Output: errors.As(err, *NotFoundError)=true
	fmt.Printf("nf.Resource=%s nf.ID=%d\n", nf.Resource, nf.ID)            // Output: nf.Resource=user nf.ID=2
}
