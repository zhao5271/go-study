package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"example.com/go-learning/internal/basics"
)

func main() {
	demoVarsZeroValues()
	demoFunctionsAndErrors()
	demoIfInit()
	demoForLoops()
	demoSwitch()
	demoSlices()
	demoNilVsEmptySliceMapAndJSON()
	demoMapIterationOrder()
}

func demoVarsZeroValues() {
	fmt.Println("\n== 1) var / := / zero values ==") // Output: == 1) var / := / zero values ==

	var count int
	var name string
	var ok bool
	fmt.Printf("count=%d name=%q ok=%v\n", count, name, ok) // Output: count=0 name="" ok=false

	// Short declaration: only inside functions.
	x := 42
	fmt.Printf("x=%d\n", x) // Output: x=42

	// Explicit type is useful when you want a specific width.
	var big int64 = 1
	fmt.Printf("big=%d\n", big) // Output: big=1

	const pi = 3.14159
	fmt.Printf("pi=%.2f\n", pi) // Output: pi=3.14
}

func demoFunctionsAndErrors() {
	fmt.Println("\n== 2) multiple returns: (value, error) ==") // Output: == 2) multiple returns: (value, error) ==

	v, err := basics.Divide(10, 2)
	fmt.Printf("Divide(10,2) => v=%d err=%v\n", v, err) // Output: Divide(10,2) => v=5 err=<nil>

	_, err = basics.Divide(10, 0)
	fmt.Printf("Divide(10,0) => err=%v\n", err) // Output: Divide(10,0) => err=divide by zero
}

func demoIfInit() {
	fmt.Println("\n== 3) if with init statement ==") // Output: == 3) if with init statement ==

	if v, err := basics.Divide(9, 3); err != nil {
		fmt.Printf("unexpected error: %v\n", err) // Output: (not printed)
	} else {
		fmt.Printf("9/3=%d\n", v) // Output: 9/3=3
	}
}

func demoForLoops() {
	fmt.Println("\n== 4) for loops (classic / while-like / range) ==") // Output: == 4) for loops (classic / while-like / range) ==

	nums := []int{10, 20, 30}

	fmt.Print("classic: ") // Output: classic:
	for i := 0; i < len(nums); i++ {
		fmt.Printf("%d ", nums[i]) // Output: 10 20 30
	}
	fmt.Println() // Output: (newline)

	fmt.Print("while-like: ") // Output: while-like:
	n := 0
	for n < 3 {
		fmt.Printf("%d ", n) // Output: 0 1 2
		n++
	}
	fmt.Println() // Output: (newline)

	fmt.Print("range(i,v): ") // Output: range(i,v):
	for i, v := range nums {
		fmt.Printf("[%d]=%d ", i, v) // Output: [0]=10 [1]=20 [2]=30
	}
	fmt.Println() // Output: (newline)
}

func demoSwitch() {
	fmt.Println("\n== 5) switch (default break + fallthrough) ==") // Output: == 5) switch (default break + fallthrough) ==

	var sb strings.Builder
	switch num := 2; num {
	case 1:
		sb.WriteString("one")
	case 2:
		sb.WriteString("two")
		fallthrough
	case 3:
		sb.WriteString("+three")
	default:
		sb.WriteString("other")
	}
	fmt.Println(sb.String()) // Output: two+three
}

func demoSlices() {
	fmt.Println("\n== 6) slices (shared backing array) ==") // Output: == 6) slices (shared backing array) ==

	s1 := []int{1, 2, 3, 4, 5}
	s2 := s1[1:3] // s2 shares backing array with s1
	s2[0] = 99
	fmt.Printf("s1=%v\n", s1) // Output: s1=[1 99 3 4 5]
	fmt.Printf("s2=%v\n", s2) // Output: s2=[99 3]

	// Avoid accidental sharing: copy into a new slice.
	s3 := make([]int, len(s2))
	copy(s3, s2)
	s3[0] = 7
	fmt.Printf("s2=%v\n", s2) // Output: s2=[99 3]
	fmt.Printf("s3=%v\n", s3) // Output: s3=[7 3]

	// Another technique: limit capacity with full slice expression to prevent appends from overwriting.
	limited := s1[1:3:3] // len=2 cap=2
	limited = append(limited, 100)
	fmt.Printf("s1=%v\n", s1)           // Output: s1=[1 99 3 4 5]
	fmt.Printf("limited=%v\n", limited) // Output: limited=[99 3 100]
}

func demoNilVsEmptySliceMapAndJSON() {
	fmt.Println("\n== 7) nil vs empty (slice/map) + JSON ==") // Output: == 7) nil vs empty (slice/map) + JSON ==

	var nilSlice []int
	emptySlice := []int{}
	fmt.Printf("len(nilSlice)=%d len(emptySlice)=%d\n", len(nilSlice), len(emptySlice)) // Output: len(nilSlice)=0 len(emptySlice)=0

	nilJSON0, _ := json.Marshal(nilSlice)
	emptyJSON, _ := json.Marshal(emptySlice)
	fmt.Printf("json(nilSlice)=%s json(emptySlice)=%s\n", nilJSON0, emptyJSON) // Output: json(nilSlice)=null json(emptySlice)=[]

	nilSlice = append(nilSlice, 1)
	nilJSON1, _ := json.Marshal(nilSlice)
	fmt.Printf("nilSlice after append=%v json=%s\n", nilSlice, nilJSON1) // Output: nilSlice after append=[1] json=[1]

	var nilMap map[string]int
	madeMap := make(map[string]int)
	fmt.Printf("nilMap==nil? %v madeMap==nil? %v\n", nilMap == nil, madeMap == nil) // Output: nilMap==nil? true madeMap==nil? false

	// nilMap["x"] = 1 // would panic: assignment to entry in nil map
	madeMap["x"] = 1
	fmt.Printf("madeMap=%v\n", madeMap) // Output: madeMap=map[x:1]
}

func demoMapIterationOrder() {
	fmt.Println("\n== 8) map iteration order ==") // Output: == 8) map iteration order ==

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Print("iter: ") // Output: iter:
	for k, v := range m {
		fmt.Printf("%s=%d ", k, v) // Output: 输出可能变化/不固定
	}
	fmt.Println() // Output: (newline)
}
