package controllers

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuthRequestValidation(t *testing.T) {
	testCases := []struct {
		AuthRequest *AuthRequest
		Valid       bool
		Error       error
	}{
		{
			AuthRequest: &AuthRequest{
				Login:    "login",
				Password: "password",
			},
			Valid: true,
		},
		{
			AuthRequest: &AuthRequest{
				Login:    "",
				Password: "password",
			},
			Valid: false,
			Error: LoginRequiredError,
		},
		{
			AuthRequest: &AuthRequest{
				Login:    "login",
				Password: "",
			},
			Valid: false,
			Error: PasswordRequiredError,
		},
	}

	for _, testCase := range testCases {
		valid, err := testCase.AuthRequest.Validate()
		require.Equal(t, testCase.Valid, valid)
		require.Equal(t, testCase.Error, err)
	}
}
