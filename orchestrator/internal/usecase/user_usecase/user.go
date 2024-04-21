package user_usecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"os"
	"time"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/dgrijalva/jwt-go"
)

func (uc UseCaseImpl) Create(ctx context.Context, request entity.RegisterRequest) error {
	exists, err := uc.userRepository.Exists(ctx, request.Login)
	if err != nil {
		return err
	}

	if exists {
		return errorz.ErrUserAlreadyExists
	}

	user := entity.User{Login: request.Login, Password: generatePasswordHash(request.Password)}
	return uc.userRepository.Create(ctx, user)
}

func (uc UseCaseImpl) Login(ctx context.Context, request entity.LoginRequest) (string, error) {
	user, err := uc.userRepository.GetUser(ctx, request.Login, generatePasswordHash(request.Password))
	if err != nil {
		return "", err
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * time.Hour).Unix(),
		},
		UserId: user.ID,
	})

	return claims.SignedString([]byte(os.Getenv("SIGNING_KEY")))
}

func (uc UseCaseImpl) ParseToken(ctx context.Context, accessToken string) (uint, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return 0, errorz.ErrInvalidSigningMethod
		}
		return []byte(os.Getenv("SIGNING_KEY")), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errorz.ErrInvalidClaimsType
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("SALT"))))
}

func (uc UseCaseImpl) ListOperations(ctx context.Context) (map[string]uint, error) {
	userID := ctx.Value("userID").(uint)
	user, err := uc.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	operations := map[string]uint{
		"+": user.Addition,
		"-": user.Subtraction,
		"*": user.Multiplication,
		"/": user.Division,
	}

	return operations, nil
}

func (uc UseCaseImpl) UpdateOperation(ctx context.Context, request entity.UpdateOperationRequest) error {
	userID := ctx.Value("userID").(uint)
	return uc.userRepository.UpdateOperation(ctx, userID, request.Operation, request.Time)
}
