package main

import (
	"Demo/day02/users"
	"fmt"
)

func main() {
	_, err := users.FindUserSentinel(2)
	fmt.Printf("sentinel err = %v\n", err)

	hidden := fmt.Errorf("hide: %v", users.ErrUserNotFound)
	fmt.Printf("hidden = %v\n", hidden)

}
