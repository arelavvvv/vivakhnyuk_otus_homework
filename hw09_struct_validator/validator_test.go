package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"
)

type UserRole string

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
				ID:    "123456789012345678901234567890123456",
				Name:  "Vladislav",
				Age:   25,
				Email: "vladislav@example.com",
				Role:  "admin",
				Phones: []string{
					"12345678901",
					"09876543210",
				},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:    "short",
				Name:  "Valeriy",
				Age:   17,
				Email: "valeriy@example",
				Role:  "user",
				Phones: []string{
					"1234567890",
				},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("length must be exactly 36")},
				{Field: "Age", Err: fmt.Errorf("must be at least 18")},
				{Field: "Email", Err: fmt.Errorf("must match regexp ^\\w+@\\w+\\.\\w+$")},
				{Field: "Role", Err: fmt.Errorf("must be one of admin,stuff")},
				{Field: "Phones", Err: fmt.Errorf("length must be exactly 11")},
			},
		},
		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1234",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: fmt.Errorf("length must be exactly 5")},
			},
		},
		{
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 201,
				Body: "Created",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("must be one of 200,404,500")},
			},
		},
		{
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
