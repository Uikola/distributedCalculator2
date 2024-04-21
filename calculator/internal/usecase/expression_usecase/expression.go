package expression_usecase

import (
	"context"
	"fmt"
	"github.com/Uikola/distributedCalculator2/calculator/pkg/polish_notation"
	"strconv"
)

func (uc UseCaseImpl) Calculate(ctx context.Context, expression string, expressionID uint) error {
	expr, err := uc.expressionRepository.GetExpressionByID(ctx, expressionID)
	if err != nil {
		return err
	}

	fmt.Println(expr.OwnerID)
	user, err := uc.userRepository.GetByID(ctx, expr.OwnerID)
	if err != nil {
		return err
	}

	operations := map[string]uint{
		"+": user.Addition,
		"-": user.Subtraction,
		"*": user.Multiplication,
		"/": user.Division,
	}

	rpn := polish_notation.ConvertToRPN(expression)
	result := polish_notation.EvalRPN(rpn, operations)

	err = uc.expressionRepository.UpdateResult(ctx, expressionID, strconv.Itoa(result))
	if err != nil {
		if unlinkErr := uc.cResourceRepository.UnlinkExpressionFromCResource(ctx, expr); unlinkErr != nil {
			return unlinkErr
		}
		if setErr := uc.expressionRepository.SetErrorStatus(ctx, expressionID); setErr != nil {
			return setErr
		}
		return err
	}

	err = uc.cResourceRepository.UnlinkExpressionFromCResource(ctx, expr)
	if err != nil {
		if setErr := uc.expressionRepository.SetErrorStatus(ctx, expressionID); setErr != nil {
			return setErr
		}
		return err
	}

	err = uc.expressionRepository.SetSuccessStatus(ctx, expressionID)
	if err != nil {
		return err
	}

	return nil
}
