package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	fmt.Println("== KP: string basics ==") // Output: == KP: string basics ==

	demoLiteralsEscapes()
	demoFmt()
	demoBuilder()
	demoCompare()
	demoMethods()
}

func demoLiteralsEscapes() {
	fmt.Println("[1] literals + escapes") // Output: [1] literals + escapes

	interp := "a\nb\tc\\\""
	fmt.Printf("interp=%q\n", interp) // Output: interp="a\nb\tc\\\""

	raw := `a\nb\tc\\\"`
	fmt.Printf("raw=%q\n", raw) // Output: raw="a\\nb\\tc\\\\\\\""

	s := "中文"
	fmt.Printf("len_bytes=%d rune_count=%d\n", len(s), utf8.RuneCountInString(s)) // Output: len_bytes=6 rune_count=2
}

func demoFmt() {
	fmt.Println("[2] fmt formatting") // Output: [2] fmt formatting

	user := "alice"
	status := 200
	msg := fmt.Sprintf("audit user=%s status=%d", user, status)
	fmt.Printf("msg=%q type=%T\n", msg, msg) // Output: msg="audit user=alice status=200" type=string

	fmt.Printf("hex=%x\n", status)    // Output: hex=c8
	fmt.Printf("pi=%8.3f\n", 3.14159) // Output: pi=   3.142
}

func demoBuilder() {
	fmt.Println("[3] strings.Builder") // Output: [3] strings.Builder

	user := "alice"
	action := "create_user"
	status := 200

	var b strings.Builder
	b.Grow(64)
	b.WriteString("audit user=")
	b.WriteString(user)
	b.WriteString(" action=")
	b.WriteString(strings.ToUpper(action))
	b.WriteString(" status=")
	b.WriteString(strconv.Itoa(status))
	fmt.Printf("%s\n", b.String()) // Output: audit user=alice action=CREATE_USER status=200

	var where strings.Builder
	var args []any
	where.WriteString("WHERE 1=1")
	where.WriteString(" AND status = ?")
	args = append(args, 1)
	fmt.Printf("sql=%s args=%v\n", where.String(), args) // Output: sql=WHERE 1=1 AND status = ? args=[1]
}

func demoCompare() {
	fmt.Println("[4] compare") // Output: [4] compare

	fmt.Printf("eq=%v\n", "admin" == "admin")                    // Output: eq=true
	fmt.Printf("fold=%v\n", strings.EqualFold("Admin", "admin")) // Output: fold=true
	fmt.Printf("cmp=%d\n", strings.Compare("a", "b"))            // Output: cmp=-1
}

func demoMethods() {
	fmt.Println("[5] common strings methods") // Output: [5] common strings methods

	q := "  admin , editor  "
	q = strings.TrimSpace(q)
	fmt.Printf("q=%q\n", q) // Output: q="admin , editor"

	fmt.Printf("hasPrefix=%v\n", strings.HasPrefix(q, "admin")) // Output: hasPrefix=true
	fmt.Printf("contains=%v\n", strings.Contains(q, "editor"))  // Output: contains=true

	parts := strings.Split(q, ",")
	fmt.Printf("parts=%v\n", parts) // Output: parts=[admin   editor]
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	fmt.Printf("trimmed=%v\n", parts) // Output: trimmed=[admin editor]

	out := strings.Join(parts, "|")
	fmt.Printf("join=%q\n", out) // Output: join="admin|editor"

	fmt.Printf("repl=%q\n", strings.ReplaceAll("a-b-c", "-", "_")) // Output: repl="a_b_c"
}
