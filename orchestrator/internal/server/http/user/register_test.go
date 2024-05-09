package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/user"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/user/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
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

func TestRegister(t *testing.T) {
	cases := []struct {
		name    string
		expCode int
		input   string

		mockErr   error
		mockInput entity.RegisterRequest

		want    map[string]string
		respErr bool
	}{
		{
			name:    "success",
			expCode: http.StatusCreated,
			input:   `{"login": "TestUser", "password": "TestPass123"}`,

			mockInput: entity.RegisterRequest{Login: "TestUser", Password: "TestPass123"},

			want: map[string]string{"response": "user created successfully"},
		},
		{
			name:    "bad json",
			expCode: http.StatusBadRequest,
			input:   `{"login": "TestUser", "password": "TestPass123"`,

			want:    map[string]string{"reason": "bad json"},
			respErr: true,
		},
		{
			name:    "bad request",
			expCode: http.StatusBadRequest,
			input:   `{"login": "testuser", "password": "testpass"}`,

			want:    map[string]string{"reason": "bad request"},
			respErr: true,
		},
		{
			name:    "user already exists",
			expCode: http.StatusConflict,
			input:   `{"login": "TestUser", "password": "TestPass123"}`,

			mockErr:   errorz.ErrUserAlreadyExists,
			mockInput: entity.RegisterRequest{Login: "TestUser", Password: "TestPass123"},

			want:    map[string]string{"reason": "user with this login already exists"},
			respErr: true,
		},
		{
			name:    "internal error",
			expCode: http.StatusInternalServerError,
			input:   `{"login": "TestUser", "password": "TestPass123"}`,

			mockErr:   errors.New("mock err"),
			mockInput: entity.RegisterRequest{Login: "TestUser", Password: "TestPass123"},

			want:    map[string]string{"reason": "internal error"},
			respErr: true,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockuserUseCase(ctrl)

			handler := user.NewHandler(mockUseCase)

			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte(tCase.input)))
			require.NoError(t, err)

			if !tCase.respErr || tCase.mockErr != nil {
				mockUseCase.EXPECT().Create(req.Context(), tCase.mockInput).Return(tCase.mockErr)
			}

			rec := httptest.NewRecorder()

			handler.Register(rec, req)
			defer rec.Result().Body.Close()

			var got map[string]string
			err = json.NewDecoder(rec.Result().Body).Decode(&got)
			require.NoError(t, err)

			require.Equal(t, tCase.expCode, rec.Result().StatusCode)
			require.Equal(t, tCase.want, got)
		})
	}
}
