package expression_usecasse

import (
	"context"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
)

//go:generate mockgen -source=usecase.go -destination=mocks/mock_expressionAndCResourceRepository.go -package=mocks expressionRepository
type expressionRepository interface {
	AddExpression(ctx context.Context, expression entity.Expression) (entity.Expression, error)
	SetErrorStatus(ctx context.Context, id uint) error
	UpdateCResource(ctx context.Context, expressionID, cResourceID uint) error
	GetExpressionByID(ctx context.Context, id uint) (entity.Expression, error)
	ListExpressions(ctx context.Context, userID uint) ([]entity.Expression, error)
	UpdateResult(ctx context.Context, expressionID uint, result string) error
	SetSuccessStatus(ctx context.Context, id uint) error
}

type cResourceRepository interface {
	AssignExpressionToCResource(ctx context.Context, expression entity.Expression) (entity.CResource, error)
	UnlinkExpressionFromCResource(ctx context.Context, expression entity.Expression) error
}

type cache interface {
	Get(ctx context.Context, key string) (string, error)
}

type UseCaseImpl struct {
	expressionRepository expressionRepository
	cResourceRepository  cResourceRepository
	cache                cache
}

func NewUseCaseImpl(expressionRepository expressionRepository, cResourceRepository cResourceRepository, cache cache) *UseCaseImpl {
	return &UseCaseImpl{
		expressionRepository: expressionRepository,
		cResourceRepository:  cResourceRepository,
		cache:                cache,
	}
}
