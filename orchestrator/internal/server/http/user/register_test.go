package user_test

import (
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/user"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateLogin(t *testing.T) {
	cases := []struct {
		name  string
		login string
		err   error
	}{
		{
			name:  "valid login",
			login: "testuser",
		},
		{
			name:  "invalid login",
			login: "t",
			err:   errorz.ErrInvalidLogin,
		},
	}

	for _, tCase := range cases {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {

			err := user.ValidateLogin(tCase.login)
			require.Equal(t, tCase.err, err)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name     string
		password string
		err      error
	}{
		{
			name:     "valid password",
			password: "TestPassword228",
		},
		{
			name:     "invalid password len",
			password: "Tp228",
			err:      errorz.ErrInvalidPassword,
		},
		{
			name:     "password without lower",
			password: "TESTPASS228",
			err:      errorz.ErrInvalidPassword,
		},
		{
			name:     "password without upper",
			password: "testpass228",
			err:      errorz.ErrInvalidPassword,
		},
		{
			name:     "password without digit",
			password: "testpass",
			err:      errorz.ErrInvalidPassword,
		},
	}

	for _, tCase := range cases {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {

			err := user.ValidatePassword(tCase.password)
			require.Equal(t, tCase.err, err)
		})
	}
}
