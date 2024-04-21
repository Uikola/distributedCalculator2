package expression_usecasse

import (
	"context"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
)

type expressionRepository interface {
	AddExpression(ctx context.Context, expression entity.Expression) (entity.Expression, error)
	SetErrorStatus(ctx context.Context, id uint) error
	UpdateCResource(ctx context.Context, expressionID, cResourceID uint) error
	GetExpressionByID(ctx context.Context, id uint) (entity.Expression, error)
	ListExpressions(ctx context.Context, userID uint) ([]entity.Expression, error)
}

type cResourceRepository interface {
	AssignExpressionToCResource(ctx context.Context, expression entity.Expression) (entity.CResource, error)
	UnlinkExpressionFromCResource(ctx context.Context, expression entity.Expression) error
}

type UseCaseImpl struct {
	expressionRepository expressionRepository
	cResourceRepository  cResourceRepository
}

func NewUseCaseImpl(expressionRepository expressionRepository, cResourceRepository cResourceRepository) *UseCaseImpl {
	return &UseCaseImpl{
		expressionRepository: expressionRepository,
		cResourceRepository:  cResourceRepository,
	}
}
