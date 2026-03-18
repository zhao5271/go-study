package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println("== Day01.5: functions + (value, error) ==") // Output: == Day01.5: functions + (value, error) ==

	v, err := divide(10, 2)
	fmt.Printf("divide(10,2) => v=%d err=%v\n", v, err) // Output: divide(10,2) => v=5 err=<nil>

	_, err = divide(10, 0)
	fmt.Printf("divide(10,0) => err=%v\n", err) // Output: divide(10,0) => err=divide by zero

	fmt.Printf("errors.Is(err, ErrDivideByZero)=%v\n", errors.Is(err, ErrDivideByZero)) // Output: errors.Is(err, ErrDivideByZero)=true
}

var ErrDivideByZero = errors.New("divide by zero")

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}
