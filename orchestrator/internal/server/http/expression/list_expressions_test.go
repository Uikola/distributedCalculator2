package expression_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/expression"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/expression/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListExpressions(t *testing.T) {
	cases := []struct {
		name    string
		expCode int

		mockErr  error
		mockResp []entity.Expression

		want      []entity.Expression
		wantIfErr map[string]string
		respErr   bool
	}{
		{
			name:    "success",
			expCode: http.StatusOK,

			mockResp: []entity.Expression{
				{ID: 1, Expression: "1 + 1", Result: "2", Status: entity.OK, CreatedAt: time.Date(2024, 5, 1, 10, 10, 10, 10, time.UTC), CalculatedAt: time.Date(2024, 5, 1, 10, 30, 10, 10, time.UTC), CalculatedBy: 1, OwnerID: 1},
				{ID: 2, Expression: "2 * 2", Result: "2", Status: entity.OK, CreatedAt: time.Date(2024, 5, 1, 10, 10, 10, 10, time.UTC), CalculatedAt: time.Date(2024, 5, 1, 10, 30, 10, 10, time.UTC), CalculatedBy: 1, OwnerID: 1},
			},

			want: []entity.Expression{
				{ID: 1, Expression: "1 + 1", Result: "2", Status: entity.OK, CreatedAt: time.Date(2024, 5, 1, 10, 10, 10, 10, time.UTC), CalculatedAt: time.Date(2024, 5, 1, 10, 30, 10, 10, time.UTC), CalculatedBy: 1, OwnerID: 1},
				{ID: 2, Expression: "2 * 2", Result: "2", Status: entity.OK, CreatedAt: time.Date(2024, 5, 1, 10, 10, 10, 10, time.UTC), CalculatedAt: time.Date(2024, 5, 1, 10, 30, 10, 10, time.UTC), CalculatedBy: 1, OwnerID: 1},
			},
		},
		{
			name:    "no expressions",
			expCode: http.StatusNoContent,

			mockErr: errorz.ErrNoExpressions,

			wantIfErr: map[string]string{"response": "you don't have any expressions"},
			respErr:   true,
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

			mockUseCase := mocks.NewMockexpressionUseCase(ctrl)

			handler := expression.NewHandler(mockUseCase)

			req, err := http.NewRequest(http.MethodGet, "/expressions", nil)
			require.NoError(t, err)

			ctx := context.WithValue(req.Context(), "userID", uint(1))

			if !tCase.respErr || tCase.mockErr != nil {
				mockUseCase.EXPECT().ListExpressions(ctx).Return(tCase.mockResp, tCase.mockErr)
			}

			rec := httptest.NewRecorder()

			handler.ListExpressions(rec, req.WithContext(ctx))
			defer rec.Result().Body.Close()

			require.Equal(t, tCase.expCode, rec.Result().StatusCode)
			if tCase.respErr {
				var got map[string]string
				err = json.NewDecoder(rec.Result().Body).Decode(&got)
				require.NoError(t, err)

				require.Equal(t, tCase.wantIfErr, got)
			} else {
				var got []entity.Expression
				err = json.NewDecoder(rec.Result().Body).Decode(&got)
				require.NoError(t, err)

				require.Equal(t, tCase.want, got)
			}
		})
	}
}
