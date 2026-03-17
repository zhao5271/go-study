package basics

import (
	"errors"
	"testing"
)

func TestDivide(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		got, err := Divide(10, 2)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if got != 5 {
			t.Fatalf("expected 5, got %d", got)
		}
	})

	t.Run("divide by zero", func(t *testing.T) {
		_, err := Divide(10, 0)
		if err == nil {
			t.Fatalf("expected error")
		}
		if !errors.Is(err, ErrDivideByZero) {
			t.Fatalf("expected ErrDivideByZero, got %v", err)
		}
	})
}
