package expression

import (
	"context"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock_expressionUsecase.go -package=mocks expressionUseCase
type expressionUseCase interface {
	AddExpression(ctx context.Context, expression string, userID uint) (uint, error)
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
