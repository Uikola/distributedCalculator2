package cresource_usecase

import (
	"context"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
)

func (uc UseCaseImpl) Create(ctx context.Context, resource entity.CResource) error {
	return uc.cResourceRepository.Create(ctx, resource)
}

func (uc UseCaseImpl) Exists(ctx context.Context, name, address string) (bool, error) {
	return uc.cResourceRepository.Exists(ctx, name, address)
}

func (uc UseCaseImpl) SetOrchestatorHealth(ctx context.Context, name string, isAlive bool) error {
	return uc.cResourceRepository.SetOrchestatorHealth(ctx, name, isAlive)
}

func (uc UseCaseImpl) Delete(ctx context.Context, name string) error {
	return uc.cResourceRepository.Delete(ctx, name)
}

func (uc UseCaseImpl) ListCResources(ctx context.Context) (map[string]string, error) {
	cResources, err := uc.cResourceRepository.ListCResources(ctx)
	if err != nil {
		return nil, err
	}

	pairs := make(map[string]string)

	for _, cResource := range cResources {
		pairs[cResource.Address] = cResource.Expression
	}

	return pairs, nil
}
