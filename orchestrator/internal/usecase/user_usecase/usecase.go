package user_usecase

import (
	"context"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/dgrijalva/jwt-go"
)

//go:generate mockgen -source=usecase.go -destination=mocks/mock_userRepository.go -package=mocks userRepository
type userRepository interface {
	Create(ctx context.Context, user entity.User) error
	Exists(ctx context.Context, login string) (bool, error)
	GetUser(ctx context.Context, login, password string) (entity.User, error)
	GetByID(ctx context.Context, id uint) (entity.User, error)
	UpdateOperation(ctx context.Context, userID uint, operation string, operationTime uint) error
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId uint `json:"userID"`
}

type UseCaseImpl struct {
	userRepository userRepository
}

func NewUseCaseImpl(userRepository userRepository) *UseCaseImpl {
	return &UseCaseImpl{userRepository: userRepository}
}
