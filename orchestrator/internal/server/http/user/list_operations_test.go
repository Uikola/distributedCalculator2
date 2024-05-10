package user_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/user"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/user/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListOperations(t *testing.T) {
	cases := []struct {
		name    string
		expCode int

		mockErr  error
		mockResp map[string]uint

		want      map[string]uint
		wantIfErr map[string]string
		respErr   bool
	}{
		{
			name:    "success",
			expCode: http.StatusOK,

			mockResp: map[string]uint{"+": 10, "-": 10, "*": 10, "/": 10},

			want: map[string]uint{"+": 10, "-": 10, "*": 10, "/": 10},
		},
		{
			name:    "internal error",
			expCode: http.StatusInternalServerError,

			mockErr: errors.New("mock err"),

			wantIfErr: map[string]string{"reason": "internal error"},
			respErr:   true,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockuserUseCase(ctrl)

			handler := user.NewHandler(mockUseCase)

			req, err := http.NewRequest(http.MethodGet, "/operations", nil)
			require.NoError(t, err)

			ctx := context.WithValue(req.Context(), "userID", uint(1))

			if !tCase.respErr || tCase.mockErr != nil {
				mockUseCase.EXPECT().ListOperations(ctx).Return(tCase.mockResp, tCase.mockErr)
			}

			rec := httptest.NewRecorder()

			handler.ListOperations(rec, req.WithContext(ctx))
			defer rec.Result().Body.Close()

			require.Equal(t, tCase.expCode, rec.Result().StatusCode)
			if tCase.respErr {
				var got map[string]string
				err = json.NewDecoder(rec.Result().Body).Decode(&got)
				require.NoError(t, err)

				require.Equal(t, tCase.wantIfErr, got)
			} else {
				var got map[string]uint
				err = json.NewDecoder(rec.Result().Body).Decode(&got)
				require.NoError(t, err)

				require.Equal(t, tCase.want, got)
			}
		})
	}
}
