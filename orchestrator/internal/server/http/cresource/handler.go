package cresource

import (
	"context"
)

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
