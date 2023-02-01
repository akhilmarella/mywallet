package utils

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestValidatePassword(t *testing.T) {
	type args struct {
		password string
	}

	type expected struct {
		valid bool
	}

	testCase := []struct {
		key string
		args
		expected
	}{
		{
			key: "if there is no uppercase in password",
			args: args{
				password: "akhil@1234",
			}, expected: expected{
				valid: false,
			},
		},
		{
			key: "if there is no lowercase in password",
			args: args{
				password: "AKHIL@1234",
			}, expected: expected{
				valid: false,
			},
		},
		{
			key: "if there is no number in password",
			args: args{
				password: "Akhil@babu",
			}, expected: expected{
				valid: false,
			},
		},
		{
			key: "if there is no symbol in password",
			args: args{
				password: "Akhil12345",
			}, expected: expected{
				valid: false,
			},
		},
		{
			key: "password is correct",
			args: args{
				password: "Akhil@1234",
			},
			expected: expected{
				valid: true,
			},
		},
		{
			key: "correct",
			args: args{
				password: "Akhilc@1245",
			},
			expected: expected{
				valid: true,
			},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.key, func(t *testing.T) {
			ok := ValidatePassword(tc.password)
			assert.Equal(t, tc.expected.valid, ok)
		})
	}
}

func TestValidateEmail(t *testing.T) {
	type args struct {
		email string
	}

	type expected struct {
		valid bool
	}

	testCase := []struct {
		key string
		args
		expected
	}{
		{key: "if @gmail is not present",
			args: args{
				email: "akhilbab123",
			}, expected: expected{
				valid: false,
			},
		}, {
			key: "email is correct ",
			args: args{
				email: "akhilbab12@gmail.com",
			}, expected: expected{
				valid: true,
			},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.key, func(t *testing.T) {
			ok := ValidateEmail(tc.email)
			assert.Equal(t, tc.expected.valid, ok)
		})
	}
}
