package expression_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/expression"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/expression/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetExpression(t *testing.T) {
	expr := entity.Expression{
		ID:           1,
		Expression:   "1 + 1",
		Status:       entity.InProgress,
		CreatedAt:    time.Date(2024, 5, 1, 10, 10, 10, 10, time.UTC),
		CalculatedBy: 1,
		OwnerID:      1,
	}

	cases := []struct {
		name    string
		expCode int

		mockErr      error
		mockResp     entity.Expression
		expressionID string

		want      entity.Expression
		wantIfErr map[string]string
		respErr   bool
	}{
		{
			name:    "success",
			expCode: http.StatusOK,

			mockResp:     expr,
			expressionID: "1",

			want: expr,
		},
		{
			name:    "invalid expression id",
			expCode: http.StatusBadRequest,

			expressionID: "test",

			wantIfErr: map[string]string{"reason": "invalid expression id"},
			respErr:   true,
		},
		{
			name:    "expression not found",
			expCode: http.StatusNotFound,

			mockErr:      errorz.ErrExpressionNotFound,
			expressionID: "1",

			wantIfErr: map[string]string{"reason": "expression not found"},
			respErr:   true,
		},
		{
			name:    "access forbidden",
			expCode: http.StatusForbidden,

			mockErr:      errorz.ErrAccessForbidden,
			expressionID: "1",

			wantIfErr: map[string]string{"reason": "access for this expression forbidden"},
			respErr:   true,
		},
		{
			name:    "internal error",
			expCode: http.StatusInternalServerError,

			mockErr:      errors.New("mock err"),
			expressionID: "1",

			wantIfErr: map[string]string{"reason": "internal error"},
			respErr:   true,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockexpressionUseCase(ctrl)

			handler := expression.NewHandler(mockUseCase)

			req, err := http.NewRequest(http.MethodGet, "/expressions/{id}", nil)
			require.NoError(t, err)

			idCtx := context.WithValue(req.Context(), "userID", uint(1))
			rCtx := chi.NewRouteContext()
			rCtx.URLParams.Add("id", tCase.expressionID)

			ctx := context.WithValue(idCtx, chi.RouteCtxKey, rCtx)

			if !tCase.respErr || tCase.mockErr != nil {
				mockUseCase.EXPECT().GetExpression(ctx, uint(1), uint(1)).Return(tCase.mockResp, tCase.mockErr)
			}

			rec := httptest.NewRecorder()

			handler.GetExpression(rec, req.WithContext(ctx))
			defer rec.Result().Body.Close()

			require.Equal(t, tCase.expCode, rec.Result().StatusCode)
			if tCase.respErr {
				var got map[string]string
				err = json.NewDecoder(rec.Result().Body).Decode(&got)
				require.NoError(t, err)

				require.Equal(t, tCase.wantIfErr, got)
			} else {
				var got entity.Expression
				err = json.NewDecoder(rec.Result().Body).Decode(&got)
				require.NoError(t, err)

				require.Equal(t, tCase.want, got)
			}
		})
	}
}
