package cresource

import (
	"context"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock_cResourceUsecase.go -package=mocks cResourceUseCase
type cResourceUseCase interface {
	ListCResources(ctx context.Context) (map[string]string, error)
}

type Handler struct {
	cResourceUseCase cResourceUseCase
}

func NewHandler(cResourceUseCase cResourceUseCase) *Handler {
	return &Handler{
		cResourceUseCase: cResourceUseCase,
	}
}
