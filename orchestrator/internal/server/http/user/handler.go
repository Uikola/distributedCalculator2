package user

import (
	"context"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
)

type userUseCase interface {
	Create(ctx context.Context, request entity.RegisterRequest) error
	Login(ctx context.Context, request entity.LoginRequest) (string, error)
	ParseToken(ctx context.Context, accessToken string) (uint, error)
	ListOperations(ctx context.Context) (map[string]uint, error)
	UpdateOperation(ctx context.Context, request entity.UpdateOperationRequest) error
}

type Handler struct {
	userUseCase userUseCase
}

func NewHandler(userUseCase userUseCase) *Handler {
	return &Handler{
		userUseCase: userUseCase,
	}
}
