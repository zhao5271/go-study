package users

import (
	"errors"
	"fmt"
)

type User struct {
	ID   int
	Name string
}

var ErrUserNotFound = errors.New("user not found")

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
	return User{}, fmt.Errorf("find user %d: %w", id, &NotFoundError{Resource: "user", ID: id})
}

func FindUserTyped(id int) (User, error) {
	if id == 1 {
		return User{ID: 1, Name: "Gopher"}, nil
	}
	return User{}, fmt.Errorf("find user %d: %w", id, &NotFoundError{Resource: "user", ID: id})
}
