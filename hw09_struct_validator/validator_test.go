package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				Age:  25,
				Role: "admin",
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrLen},
				{Field: "Email", Err: ErrRegexp},
			},
		},
		{
			in: User{
				Age:   15,
				Email: "test@gmail.com",
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrLen},
				{Field: "Age", Err: ErrMinLimit},
				{Field: "Role", Err: ErrInSet},
			},
		},
		{
			in: User{
				Age: 60,
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrLen},
				{Field: "Age", Err: ErrMaxLimit},
				{Field: "Email", Err: ErrRegexp},
				{Field: "Role", Err: ErrInSet},
			},
		},
		{
			in:          Response{Code: 200},
			expectedErr: nil,
		},
		{
			in: Response{Code: 401},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: ErrInSet},
			},
		},
		{
			in: User{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Age:   18,
				Role:  "admin",
				Email: "test@gmail.com",
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:   "550e8400-e29b-41d4-a716-446655440000",
				Age:  18,
				Role: "manager",
			},
			expectedErr: ValidationErrors{
				{Field: "Email", Err: ErrRegexp},
				{Field: "Role", Err: ErrInSet},
			},
		},
		{
			in: User{
				ID:   "550e8400-e29b-41d4-a716-446655440000",
				Age:  18,
				Role: "admin",
				Phones: []string{
					"79163272727", "79993272727",
				},
				Email: "test@gmail.com",
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:   "550e8400-e29b-41d4-a716-446655440000",
				Age:  18,
				Role: "admin",
				Phones: []string{
					"79163272727", "7999327272",
				},
				Email: "test@gmail.com",
			},
			expectedErr: ValidationErrors{
				{Field: "Phones[1]", Err: ErrLen},
			},
		},
		{
			in:          App{Version: "1.0.1"},
			expectedErr: nil,
		},
		{
			in: App{Version: "1.0.1-beta"},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: ErrLen},
			},
		},
		{
			in:          123,
			expectedErr: ErrNotStruct,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				if err != nil {
					t.Errorf("got unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Errorf("expected error %v, but got nil", tt.expectedErr)
				return
			}

			var gotVE, wantVE ValidationErrors
			if errors.As(tt.expectedErr, &wantVE) {
				if errors.As(err, &gotVE) {
					compareErrors(t, gotVE, wantVE)
					return
				}
				t.Fatalf("expected ValidationErrors, but got: %T", err)
			}

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error to be %v, but got %v", tt.expectedErr, err)
			}
		})
	}
}

func compareErrors(t *testing.T, got, want ValidationErrors) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("expected %d errors, got %d. Full error: %v", len(want), len(got), got)
		return
	}

	for _, w := range want {
		found := false
		for _, g := range got {
			if g.Field == w.Field && errors.Is(g.Err, w.Err) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("not found expected error for field %s: %v", w.Field, w.Err)
		}
	}
}
