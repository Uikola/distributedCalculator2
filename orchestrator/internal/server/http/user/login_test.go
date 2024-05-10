package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/user"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/user/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLogin(t *testing.T) {
	cases := []struct {
		name    string
		expCode int
		input   string

		mockErr   error
		mockInput entity.LoginRequest
		mockResp  string

		want    map[string]string
		respErr bool
	}{
		{
			name:    "success",
			expCode: http.StatusOK,
			input:   `{"login": "TestUser", "password": "TestPass123"}`,

			mockInput: entity.LoginRequest{Login: "TestUser", Password: "TestPass123"},
			mockResp:  "megatestToken",

			want: map[string]string{"token": "megatestToken"},
		},
		{
			name:    "bad json",
			expCode: http.StatusBadRequest,
			input:   `{"login": "TestUser", "password": "TestPass123"`,

			want:    map[string]string{"reason": "bad json"},
			respErr: true,
		},
		{
			name:    "user not found",
			expCode: http.StatusNotFound,
			input:   `{"login": "TestUser", "password": "TestPass123"}`,

			mockInput: entity.LoginRequest{Login: "TestUser", Password: "TestPass123"},
			mockErr:   errorz.ErrUserNotFound,

			want:    map[string]string{"reason": "user not found"},
			respErr: true,
		},
		{
			name:    "internal error",
			expCode: http.StatusInternalServerError,
			input:   `{"login": "TestUser", "password": "TestPass123"}`,

			mockInput: entity.LoginRequest{Login: "TestUser", Password: "TestPass123"},
			mockErr:   errors.New("mock err"),

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

			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(tCase.input)))
			require.NoError(t, err)

			if !tCase.respErr || tCase.mockErr != nil {
				mockUseCase.EXPECT().Login(req.Context(), tCase.mockInput).Return(tCase.mockResp, tCase.mockErr)
			}

			rec := httptest.NewRecorder()

			handler.Login(rec, req)
			defer rec.Result().Body.Close()

			var got map[string]string
			err = json.NewDecoder(rec.Result().Body).Decode(&got)
			require.NoError(t, err)

			require.Equal(t, tCase.expCode, rec.Result().StatusCode)
			require.Equal(t, tCase.want, got)
		})
	}
}
