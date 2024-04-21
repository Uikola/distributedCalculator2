package expression

import (
	"context"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
)

type expressionUseCase interface {
	AddExpression(ctx context.Context, expression entity.Expression) (uint, error)
	GetExpression(ctx context.Context, userID, expressionID uint) (entity.Expression, error)
	ListExpressions(ctx context.Context) ([]entity.Expression, error)
	GetResult(ctx context.Context, expressionID uint) (string, error)
}

type Handler struct {
	expressionUseCase expressionUseCase
}

func NewHandler(expressionUseCase expressionUseCase) *Handler {
	return &Handler{
		expressionUseCase: expressionUseCase,
	}
}
