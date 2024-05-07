package expression_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/expression"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/expression/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateExpression(t *testing.T) {
	cases := []struct {
		name string
		expr string
		err  error
	}{
		{
			name: "valid expression",
			expr: "2 + 2",
		},
		{
			name: "invalid expression syntax",
			expr: "1 + 2 + test",
			err:  errorz.ErrInvalidExpression,
		},
		{
			name: "invalid expression",
			expr: "1 + 2()3(",
			err:  errorz.ErrInvalidExpression,
		},
	}

	for _, tCase := range cases {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {

			err := expression.ValidateExpression(tCase.expr)
			require.Equal(t, tCase.err, err)
		})
	}
}

func TestAddExpression(t *testing.T) {

	cases := []struct {
		name    string
		input   string
		expCode int

		mockErr     error
		useCaseResp uint

		want    string
		errWant map[string]string
		respErr bool
	}{
		{
			name:    "success",
			input:   `{"expression": "1 + 1"}`,
			expCode: http.StatusOK,

			useCaseResp: 1,

			want: "1",
		},
		{
			name:    "bad json",
			input:   `{"expression": "1 + 1"`,
			expCode: http.StatusBadRequest,

			errWant: map[string]string{"reason": "bad json"},
			respErr: true,
		},
		{
			name:    "invalid expression",
			input:   `{"expression": "test"}`,
			expCode: http.StatusBadRequest,

			errWant: map[string]string{"reason": "bad request"},
			respErr: true,
		},
		{
			name:    "no avalible computing resources",
			input:   `{"expression": "1 + 1"}`,
			expCode: http.StatusNoContent,

			mockErr: errorz.ErrNoAvailableResources,

			errWant: map[string]string{"reason": "no available computing resources"},
			respErr: true,
		},
		{
			name:    "use case error",
			input:   `{"expression": "1 + 1"}`,
			expCode: http.StatusInternalServerError,

			mockErr: errors.New("mock err"),

			errWant: map[string]string{"reason": "internal error"},
			respErr: true,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockexpressionUseCase(ctrl)

			handler := expression.NewHandler(mockUseCase)

			req, err := http.NewRequest(http.MethodPost, "/api/calculate", bytes.NewBuffer([]byte(tCase.input)))
			require.NoError(t, err)

			ctx := context.WithValue(req.Context(), "userID", uint(1))

			if !tCase.respErr || tCase.mockErr != nil {
				mockUseCase.EXPECT().AddExpression(ctx, "1 + 1", uint(1)).Return(tCase.useCaseResp, tCase.mockErr)
			}

			rec := httptest.NewRecorder()

			handler.AddExpression(rec, req.WithContext(ctx))
			defer rec.Result().Body.Close()

			var got map[string]string
			err = json.NewDecoder(rec.Result().Body).Decode(&got)
			require.NoError(t, err)

			require.Equal(t, tCase.expCode, rec.Result().StatusCode)
			if tCase.respErr {
				require.Equal(t, tCase.errWant, got)
			} else {
				require.Equal(t, tCase.want, got["expression_id"])
			}
		})
	}
}
