package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := s1[1:3]
	s2[0] = 99
	fmt.Printf("s1:%v,s2:%v\n", s1, s2)

	s3 := make([]int, len(s2))
	copy(s3, s2)
	s3[0] = 7
	fmt.Printf("s2=%v\n", s2)
	fmt.Printf("s3=%v\n", s3)

	var nilSlice []int
	emptySlice := []int{}
	nilJSON, _ := json.Marshal(nilSlice)
	emptyJSON, _ := json.Marshal(emptySlice)
	fmt.Printf("json(nilSlice)=%s,json(nilJSON)%s\n", nilJSON, emptyJSON)

	var nilMap map[string]int
	m := make(map[string]int)
	fmt.Printf("nilMap==nil? %v m==nil? %v\n", nilMap == nil, m == nil) // Output: nilMap==nil? true m==nil? false

}
