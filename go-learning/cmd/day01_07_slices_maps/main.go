package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	fmt.Println("== Day01.7: slices/maps + nil/empty + map order ==") // Output: == Day01.7: slices/maps + nil/empty + map order ==

	// Slice shares backing array
	s1 := []int{1, 2, 3, 4, 5}
	s2 := s1[1:3]
	s2[0] = 99
	fmt.Printf("s1=%v\n", s1) // Output: s1=[1 99 3 4 5]
	fmt.Printf("s2=%v\n", s2) // Output: s2=[99 3]

	// Avoid accidental sharing: copy
	s3 := make([]int, len(s2))
	copy(s3, s2)
	s3[0] = 7
	fmt.Printf("s2=%v\n", s2) // Output: s2=[99 3]
	fmt.Printf("s3=%v\n", s3) // Output: s3=[7 3]

	// nil vs empty slice in JSON
	var nilSlice []int
	emptySlice := []int{}
	nilJSON, _ := json.Marshal(nilSlice)
	emptyJSON, _ := json.Marshal(emptySlice)
	fmt.Printf("json(nilSlice)=%s json(emptySlice)=%s\n", nilJSON, emptyJSON) // Output: json(nilSlice)=null json(emptySlice)=[]

	// nil map cannot be assigned into; use make
	var nilMap map[string]int
	m := make(map[string]int)
	fmt.Printf("nilMap==nil? %v m==nil? %v\n", nilMap == nil, m == nil) // Output: nilMap==nil? true m==nil? false
	// nilMap["x"] = 1 // would panic: assignment to entry in nil map
	m["x"] = 1
	fmt.Printf("m=%v\n", m) // Output: m=map[x:1]

	// map iteration order is not deterministic
	m2 := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Print("iter: ") // Output: iter:
	for k, v := range m2 {
		fmt.Printf("%s=%d ", k, v) // Output: 输出可能变化/不固定
	}
	fmt.Println() // Output: (newline)
}
