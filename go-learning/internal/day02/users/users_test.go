package users

import (
	"errors"
	"testing"
)

func TestFindUserSentinel(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		wantUserID int
		wantErrIs  error
	}{
		{
			name:       "found",
			id:         1,
			wantUserID: 1,
			wantErrIs:  nil,
		},
		{
			name:       "not found",
			id:         2,
			wantUserID: 0,
			wantErrIs:  ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := FindUserSentinel(tt.id)
			if user.ID != tt.wantUserID {
				t.Fatalf("user.ID=%d, want %d", user.ID, tt.wantUserID)
			}

			if tt.wantErrIs == nil {
				if err != nil {
					t.Fatalf("err=%v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("err=nil, want non-nil")
			}
			if !errors.Is(err, tt.wantErrIs) {
				t.Fatalf("errors.Is(err, %v)=false, err=%v", tt.wantErrIs, err)
			}
		})
	}
}

func TestFindUserTyped(t *testing.T) {
	_, err := FindUserTyped(2)
	if err == nil {
		t.Fatalf("err=nil, want non-nil")
	}

	var nf *NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("errors.As(err, *NotFoundError)=false, err=%v", err)
	}
	if nf.Resource != "user" || nf.ID != 2 {
		t.Fatalf("nf=%+v, want Resource=user ID=2", nf)
	}
}
