package users

import (
	"errors"
	"fmt"
)

type User struct {
	ID   int
	Name string
}

// ErrUserNotFound is a sentinel error. Callers should match it with errors.Is
// (not ==), because we return it wrapped with context.
var ErrUserNotFound = errors.New("user not found")

// NotFoundError is a typed error. Callers can inspect it with errors.As.
type NotFoundError struct {
	Resource string
	ID       int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s %d not found", e.Resource, e.ID)
}

func FindUserSentinel(id int) (User, error) {
	if id == 1 {
		return User{ID: 1, Name: "Gopher"}, nil
	}
	return User{}, fmt.Errorf("find user %d: %w", id, ErrUserNotFound)
}

func FindUserTyped(id int) (User, error) {
	if id == 1 {
		return User{ID: 1, Name: "Gopher"}, nil
	}
	return User{}, fmt.Errorf("find user %d: %w", id, &NotFoundError{Resource: "user", ID: id})
}
