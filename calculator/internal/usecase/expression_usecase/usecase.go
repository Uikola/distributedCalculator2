package expression_usecase

import (
	"context"
	"github.com/Uikola/distributedCalculator2/calculator/internal/entity"
)

type expressionRepository interface {
	UpdateResult(ctx context.Context, expressionID uint, result string) error
	GetExpressionByID(ctx context.Context, id uint) (entity.Expression, error)
	SetErrorStatus(ctx context.Context, id uint) error
	SetSuccessStatus(ctx context.Context, id uint) error
}

type cResourceRepository interface {
	UnlinkExpressionFromCResource(ctx context.Context, expression entity.Expression) error
}

type userRepository interface {
	GetByID(ctx context.Context, id uint) (entity.User, error)
}

type UseCaseImpl struct {
	expressionRepository expressionRepository
	cResourceRepository  cResourceRepository
	userRepository       userRepository
}

func NewUseCaseImpl(expressionRepository expressionRepository, cResourceRepository cResourceRepository, userRepository userRepository) *UseCaseImpl {
	return &UseCaseImpl{expressionRepository: expressionRepository, cResourceRepository: cResourceRepository, userRepository: userRepository}
}
