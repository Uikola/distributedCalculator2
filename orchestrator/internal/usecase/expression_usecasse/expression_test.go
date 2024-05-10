package expression_usecasse_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/usecase/expression_usecasse"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/usecase/expression_usecasse/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetExpression(t *testing.T) {
	cases := []struct {
		name string

		expressionID uint
		userID       uint
		mockErr      error
		mockResp     entity.Expression

		want    entity.Expression
		wantErr error
	}{
		{
			name: "success",

			expressionID: 1,
			userID:       1,
			mockResp:     entity.Expression{ID: 1, Expression: "1 + 1", Status: entity.InProgress, OwnerID: 1},

			want: entity.Expression{ID: 1, Expression: "1 + 1", Status: entity.InProgress, OwnerID: 1},
		},
		{
			name: "mock error",

			expressionID: 1,
			userID:       1,
			mockErr:      errors.New("mock err"),

			wantErr: errors.New("mock err"),
		},
		{
			name: "access forbidden",

			expressionID: 1,
			userID:       1,
			mockResp:     entity.Expression{ID: 1, Expression: "1 + 1", Status: entity.InProgress, OwnerID: 2},

			wantErr: errorz.ErrAccessForbidden,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExpressionRepository := mocks.NewMockexpressionRepository(ctrl)
			mockCResourceRepository := mocks.NewMockcResourceRepository(ctrl)

			mockExpressionRepository.EXPECT().GetExpressionByID(context.Background(), tCase.expressionID).Return(tCase.mockResp, tCase.mockErr)

			usecase := expression_usecasse.NewUseCaseImpl(mockExpressionRepository, mockCResourceRepository)

			got, gotErr := usecase.GetExpression(context.Background(), tCase.userID, tCase.expressionID)

			require.Equal(t, tCase.want, got)
			require.Equal(t, tCase.wantErr, gotErr)
		})
	}
}

func TestListExpressions(t *testing.T) {
	cases := []struct {
		name string

		userID   uint
		mockErr  error
		mockResp []entity.Expression

		want    []entity.Expression
		wantErr error
	}{
		{
			name: "success",

			userID: 1,
			mockResp: []entity.Expression{
				{ID: 1, Expression: "1 + 1", Result: "2", Status: entity.OK, CalculatedBy: 1, OwnerID: 1},
				{ID: 2, Expression: "2 + 2", Status: entity.InProgress, CalculatedBy: 2, OwnerID: 1},
			},

			want: []entity.Expression{
				{ID: 1, Expression: "1 + 1", Result: "2", Status: entity.OK, CalculatedBy: 1, OwnerID: 1},
				{ID: 2, Expression: "2 + 2", Status: entity.InProgress, CalculatedBy: 2, OwnerID: 1},
			},
		},
		{
			name: "mock error",

			userID:  1,
			mockErr: errors.New("mock err"),

			want:    []entity.Expression{},
			wantErr: errors.New("mock err"),
		},
		{
			name: "no expressions",

			userID: 1,

			want:    []entity.Expression{},
			wantErr: errorz.ErrNoExpressions,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExpressionRepository := mocks.NewMockexpressionRepository(ctrl)
			mockCResourceRepository := mocks.NewMockcResourceRepository(ctrl)

			mockExpressionRepository.EXPECT().ListExpressions(context.WithValue(context.Background(), "userID", tCase.userID), tCase.userID).Return(tCase.mockResp, tCase.mockErr)

			usecase := expression_usecasse.NewUseCaseImpl(mockExpressionRepository, mockCResourceRepository)

			got, gotErr := usecase.ListExpressions(context.WithValue(context.Background(), "userID", tCase.userID))

			require.Equal(t, tCase.want, got)
			require.Equal(t, tCase.wantErr, gotErr)
		})
	}
}

func TestGetResult(t *testing.T) {
	cases := []struct {
		name string

		expressionID uint
		userID       uint
		mockErr      error
		mockResp     entity.Expression

		want    string
		wantErr error
	}{
		{
			name: "success",

			expressionID: 1,
			userID:       1,
			mockResp:     entity.Expression{ID: 1, Expression: "1 + 1", Result: "2", Status: entity.OK, CalculatedBy: 1, OwnerID: 1},

			want: "2",
		},
		{
			name: "mock err",

			expressionID: 1,
			userID:       1,
			mockErr:      errors.New("mock err"),

			wantErr: errors.New("mock err"),
		},
		{
			name: "access forbidden",

			expressionID: 1,
			userID:       1,
			mockResp:     entity.Expression{ID: 1, Expression: "1 + 1", Result: "2", Status: entity.OK, CalculatedBy: 1, OwnerID: 2},

			wantErr: errorz.ErrAccessForbidden,
		},
		{
			name: "calculation in progress",

			expressionID: 1,
			userID:       1,
			mockResp:     entity.Expression{ID: 1, Expression: "1 + 1", Result: "", Status: entity.InProgress, CalculatedBy: 1, OwnerID: 1},

			wantErr: errorz.ErrEvaluationInProgress,
		},
		{
			name: "calculation error",

			expressionID: 1,
			userID:       1,
			mockResp:     entity.Expression{ID: 1, Expression: "1 + 1", Result: "", Status: entity.Error, CalculatedBy: 1, OwnerID: 1},

			wantErr: errorz.ErrEvaluation,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExpressionRepository := mocks.NewMockexpressionRepository(ctrl)
			mockCResourceRepository := mocks.NewMockcResourceRepository(ctrl)

			mockExpressionRepository.EXPECT().GetExpressionByID(context.WithValue(context.Background(), "userID", tCase.userID), tCase.expressionID).Return(tCase.mockResp, tCase.mockErr)

			usecase := expression_usecasse.NewUseCaseImpl(mockExpressionRepository, mockCResourceRepository)

			got, gotErr := usecase.GetResult(context.WithValue(context.Background(), "userID", tCase.userID), tCase.expressionID)

			require.Equal(t, tCase.want, got)
			require.Equal(t, tCase.wantErr, gotErr)
		})
	}
}
