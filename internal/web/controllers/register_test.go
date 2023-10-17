package controllers

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegisterRequestValidation(t *testing.T) {
	testCases := []struct {
		Request *RegisterRequest
		Valid   bool
		Error   error
	}{
		{
			Request: &RegisterRequest{
				Login:           "login",
				Password:        "password",
				ConfirmPassword: "password",
			},
			Valid: true,
		},
		{
			Request: &RegisterRequest{
				Login:           "",
				Password:        "password",
				ConfirmPassword: "password",
			},
			Valid: false,
			Error: LoginRequiredError,
		},
		{
			Request: &RegisterRequest{
				Login:           "login",
				Password:        "",
				ConfirmPassword: "password",
			},
			Valid: false,
			Error: PasswordRequiredError,
		},
		{
			Request: &RegisterRequest{
				Login:           "login",
				Password:        "password",
				ConfirmPassword: "",
			},
			Valid: false,
			Error: PasswordRequiredError,
		},
		{
			Request: &RegisterRequest{
				Login:           "login",
				Password:        "password",
				ConfirmPassword: "password1",
			},
			Valid: false,
			Error: PasswordNotEqualError,
		},
	}

	for _, testCase := range testCases {
		valid, err := testCase.Request.Validate()
		require.Equal(t, testCase.Valid, valid)
		require.Equal(t, testCase.Error, err)
	}
}

func TestRegisterControllerCreation(t *testing.T) {
	c := NewRegisterController(nil, nil)
	require.Equal(t, c, &RegisterController{})
}

func TestRegisterControllerGroup(t *testing.T) {
	c := NewRegisterController(nil, nil)
	require.Equal(t, "/auth", c.GetGroup())
}
