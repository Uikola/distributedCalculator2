package user_test

import (
	"bytes"
	"context"
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

func TestValidateUpdateOperation(t *testing.T) {
	cases := []struct {
		name string
		req  entity.UpdateOperationRequest
		err  error
	}{
		{
			name: "valid operation and time",
			req:  entity.UpdateOperationRequest{Operation: "+", Time: 10},
		},
		{
			name: "invalid operation",
			req:  entity.UpdateOperationRequest{Operation: "test", Time: 10},
			err:  errorz.ErrInvalidOperation,
		},
	}

	for _, tCase := range cases {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {

			err := user.ValidateUpdateOperation(tCase.req)
			require.Equal(t, tCase.err, err)
		})
	}
}

func TestUpdateOperation(t *testing.T) {
	cases := []struct {
		name    string
		expCode int
		input   string

		mockErr   error
		mockInput entity.UpdateOperationRequest

		want    map[string]string
		respErr bool
	}{
		{
			name:    "success",
			expCode: http.StatusOK,
			input:   `{"operation": "+", "time": 20}`,

			mockInput: entity.UpdateOperationRequest{Operation: "addition", Time: 20},

			want: map[string]string{"response": "operation updated successfully"},
		},
		{
			name:    "bad json",
			expCode: http.StatusBadRequest,
			input:   `{"operation": "+", "time": 20`,

			want:    map[string]string{"reason": "bad json"},
			respErr: true,
		},
		{
			name:    "bad request",
			expCode: http.StatusBadRequest,
			input:   `{"operation": "test", "time": 20}`,

			want:    map[string]string{"reason": "bad request"},
			respErr: true,
		},
		{
			name:    "internal error",
			expCode: http.StatusInternalServerError,
			input:   `{"operation": "+", "time": 20}`,

			mockErr:   errors.New("mock err"),
			mockInput: entity.UpdateOperationRequest{Operation: "addition", Time: 20},

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

			req, err := http.NewRequest(http.MethodPut, "/operations", bytes.NewBuffer([]byte(tCase.input)))
			require.NoError(t, err)

			ctx := context.WithValue(req.Context(), "userID", uint(1))

			if !tCase.respErr || tCase.mockErr != nil {
				mockUseCase.EXPECT().UpdateOperation(ctx, tCase.mockInput).Return(tCase.mockErr)
			}

			rec := httptest.NewRecorder()

			handler.UpdateOperation(rec, req.WithContext(ctx))
			defer rec.Result().Body.Close()

			var got map[string]string
			err = json.NewDecoder(rec.Result().Body).Decode(&got)
			require.NoError(t, err)

			require.Equal(t, tCase.expCode, rec.Result().StatusCode)
			require.Equal(t, tCase.want, got)
		})
	}
}
