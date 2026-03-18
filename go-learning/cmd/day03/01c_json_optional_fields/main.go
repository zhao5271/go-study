package main

import (
	"encoding/json"
	"fmt"
)

// 在 HTTP API 里，为了区分“字段缺失” vs “字段显式给了零值”，常用指针字段表示 optional。
type PatchUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Admin *bool   `json:"admin,omitempty"`
}

func main() {
	fmt.Println("== Day03.1c: JSON optional fields (pointer) ==") // Output: == Day03.1c: JSON optional fields (pointer) ==

	rawMissing := []byte(`{}`)
	var r1 PatchUserRequest
	_ = json.Unmarshal(rawMissing, &r1)
	fmt.Printf("missing: nameNil=%v adminNil=%v\n", r1.Name == nil, r1.Admin == nil) // Output: missing: nameNil=true adminNil=true

	rawZero := []byte(`{"admin": false}`)
	var r2 PatchUserRequest
	_ = json.Unmarshal(rawZero, &r2)
	fmt.Printf("explicit false: adminNil=%v adminVal=%v\n", r2.Admin == nil, *r2.Admin) // Output: explicit false: adminNil=false adminVal=false

	name := "Alice"
	admin := true
	out, _ := json.Marshal(PatchUserRequest{Name: &name, Admin: &admin})
	fmt.Printf("marshal: %s\n", string(out)) // Output: marshal: {"name":"Alice","admin":true}
}

