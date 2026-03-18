package main

import "fmt"

type User struct {
	ID    int
	Name  string
	Admin bool
}

func main() {
	fmt.Println("== Day03.1a: struct + zero value ==") // Output: == Day03.1a: struct + zero value ==

	var u User
	fmt.Printf("u=%+v\n", u) // Output: u={ID:0 Name: Admin:false}

	u2 := User{ID: 1, Name: "Alice", Admin: true}
	fmt.Printf("u2=%+v\n", u2) // Output: u2={ID:1 Name:Alice Admin:true}

	u3 := User{Name: "Bob"} // 没赋值的字段会是零值
	fmt.Printf("u3=%+v\n", u3) // Output: u3={ID:0 Name:Bob Admin:false}
}

