package cresource_usecase

import (
	"context"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
)

type cResourceRepository interface {
	Create(ctx context.Context, resource entity.CResource) error
	Exists(ctx context.Context, name, address string) (bool, error)
	SetOrchestatorHealth(ctx context.Context, name string, isAlive bool) error
	Delete(ctx context.Context, name string) error
	ListCResources(ctx context.Context) ([]entity.CResource, error)
}

type UseCaseImpl struct {
	cResourceRepository cResourceRepository
}

func NewUseCaseImpl(cResourceRepository cResourceRepository) *UseCaseImpl {
	return &UseCaseImpl{cResourceRepository: cResourceRepository}
}
