package expression_usecasse

import (
	"context"
	"time"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/Uikola/distributedCalculator2/orchestrator/pkg/grpc/client/expression"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (uc UseCaseImpl) AddExpression(ctx context.Context, exp string, userID uint) (uint, error) {
	expr := entity.Expression{
		Expression: exp,
		Status:     entity.InProgress,
		CreatedAt:  time.Now(),
		OwnerID:    userID,
	}

	expr, err := uc.expressionRepository.AddExpression(ctx, expr)
	if err != nil {
		return 0, err
	}

	cResource, err := uc.cResourceRepository.AssignExpressionToCResource(ctx, expr)
	if err != nil {
		if setErr := uc.expressionRepository.SetErrorStatus(ctx, expr.ID); setErr != nil {
			return 0, setErr
		}
		return 0, err
	}

	result, err := uc.cache.Get(ctx, expr.Expression)
	if err == nil {
		if updateErr := uc.expressionRepository.UpdateResult(ctx, expr.ID, result); updateErr != nil {
			return 0, updateErr
		}
		if setErr := uc.expressionRepository.SetSuccessStatus(ctx, expr.ID); setErr != nil {
			return 0, setErr
		}
		return expr.ID, nil
	}

	err = uc.expressionRepository.UpdateCResource(ctx, expr.ID, cResource.ID)
	if err != nil {
		if unlinkErr := uc.cResourceRepository.UnlinkExpressionFromCResource(ctx, expr); unlinkErr != nil {
			return 0, unlinkErr
		}
		if setErr := uc.expressionRepository.SetErrorStatus(ctx, expr.ID); setErr != nil {
			return 0, setErr
		}
		return 0, err
	}

	calcConn, err := grpc.NewClient(cResource.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		if unlinkErr := uc.cResourceRepository.UnlinkExpressionFromCResource(ctx, expr); unlinkErr != nil {
			return 0, unlinkErr
		}
		if setErr := uc.expressionRepository.SetErrorStatus(ctx, expr.ID); setErr != nil {
			return 0, setErr
		}
		return 0, err
	}

	expressionClient := expression.NewClient(calcConn, uc.expressionRepository, uc.cResourceRepository)
	expressionClient.Calculate(ctx, expr, calcConn)

	return expr.ID, nil
}

func (uc UseCaseImpl) GetExpression(ctx context.Context, userID, expressionID uint) (entity.Expression, error) {
	expr, err := uc.expressionRepository.GetExpressionByID(ctx, expressionID)
	if err != nil {
		return entity.Expression{}, err
	}

	if userID != expr.OwnerID {
		return entity.Expression{}, errorz.ErrAccessForbidden
	}

	return expr, nil
}

func (uc UseCaseImpl) ListExpressions(ctx context.Context) ([]entity.Expression, error) {
	userID := ctx.Value("userID").(uint)

	expressions, err := uc.expressionRepository.ListExpressions(ctx, userID)
	if err != nil {
		return []entity.Expression{}, err
	}

	if len(expressions) == 0 {
		return []entity.Expression{}, errorz.ErrNoExpressions
	}

	return expressions, nil
}

func (uc UseCaseImpl) GetResult(ctx context.Context, expressionID uint) (string, error) {
	userID := ctx.Value("userID").(uint)

	expr, err := uc.expressionRepository.GetExpressionByID(ctx, expressionID)
	if err != nil {
		return "", err
	}

	if userID != expr.OwnerID {
		return "", errorz.ErrAccessForbidden
	}

	switch expr.Status {
	case entity.InProgress:
		return "", errorz.ErrEvaluationInProgress
	case entity.Error:
		return "", errorz.ErrEvaluation
	}

	return expr.Result, nil
}
