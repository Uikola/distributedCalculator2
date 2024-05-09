package expression_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/expression"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/expression/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetResult(t *testing.T) {
	cases := []struct {
		name    string
		expCode int

		mockErr      error
		mockResp     string
		expressionID string

		want    map[string]string
		respErr bool
	}{
		{
			name:    "success",
			expCode: 200,

			mockResp:     "2",
			expressionID: "1",

			want: map[string]string{"result": "2"},
		},
		{
			name:    "invalid expression id",
			expCode: http.StatusBadRequest,

			expressionID: "test",

			want:    map[string]string{"reason": "invalid expression id"},
			respErr: true,
		},
		{
			name:    "expression not found",
			expCode: http.StatusNotFound,

			mockErr:      errorz.ErrExpressionNotFound,
			expressionID: "1",

			want:    map[string]string{"reason": "expression not found"},
			respErr: true,
		},
		{
			name:    "access forbidden",
			expCode: http.StatusForbidden,

			mockErr:      errorz.ErrAccessForbidden,
			expressionID: "1",

			want:    map[string]string{"reason": "access for this expression forbidden"},
			respErr: true,
		},
		{
			name:    "evaluation in progress",
			expCode: http.StatusOK,

			mockErr:      errorz.ErrEvaluationInProgress,
			expressionID: "1",

			want: map[string]string{"response": "your expression is being calculated, wait a bit"},
		},
		{
			name:    "evaluation error",
			expCode: http.StatusOK,

			mockErr:      errorz.ErrEvaluation,
			expressionID: "1",

			want:    map[string]string{"response": "an error occurred during the calculation of your expression, try again"},
			respErr: true,
		},
		{
			name:    "internal error",
			expCode: http.StatusInternalServerError,

			mockErr:      errors.New("mock err"),
			expressionID: "1",

			want:    map[string]string{"reason": "internal error"},
			respErr: true,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockexpressionUseCase(ctrl)

			handler := expression.NewHandler(mockUseCase)

			req, err := http.NewRequest(http.MethodGet, "/results/{id}", nil)
			require.NoError(t, err)

			idCtx := context.WithValue(req.Context(), "userID", uint(1))
			rCtx := chi.NewRouteContext()
			rCtx.URLParams.Add("id", tCase.expressionID)

			ctx := context.WithValue(idCtx, chi.RouteCtxKey, rCtx)

			if !tCase.respErr || tCase.mockErr != nil {
				mockUseCase.EXPECT().GetResult(ctx, uint(1)).Return(tCase.mockResp, tCase.mockErr)
			}

			rec := httptest.NewRecorder()

			handler.GetResult(rec, req.WithContext(ctx))
			defer rec.Result().Body.Close()

			require.Equal(t, tCase.expCode, rec.Result().StatusCode)

			var got map[string]string
			err = json.NewDecoder(rec.Result().Body).Decode(&got)
			require.NoError(t, err)

			require.Equal(t, tCase.want, got)
		})
	}
}
