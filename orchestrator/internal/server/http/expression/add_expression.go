package expression

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/rs/zerolog/log"
	"net/http"
	"regexp"
)

func (h Handler) AddExpression(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var addExpressionRequest entity.AddExpressionRequest
	if err := json.NewDecoder(r.Body).Decode(&addExpressionRequest); err != nil {
		log.Error().Err(err).Msg("failed to parse request")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad json"})
		return
	}

	if err := ValidateExpression(addExpressionRequest.Expression); err != nil {
		log.Error().Err(err).Msg("failed to parse validate the expression")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad request"})
		return
	}

	userID := ctx.Value("userID").(uint)

	id, err := h.expressionUseCase.AddExpression(ctx, addExpressionRequest.Expression, userID)
	switch {
	case errors.Is(errorz.ErrNoAvailableResources, err):
		log.Error().Err(err).Msg("no available comptuing resources")
		w.WriteHeader(http.StatusNoContent)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "no available computing resources"})
		return
	case err != nil:
		log.Error().Err(err).Msg("failed to add expression")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"expression_id": fmt.Sprintf("%d", id), "msg": "use this id to find out the result"})
}

func ValidateExpression(expr string) error {
	re := regexp.MustCompile(`[^0-9+\-*/() ]`)

	if re.MatchString(expr) {
		return errorz.ErrInvalidExpression
	}

	expression, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return errorz.ErrInvalidExpression
	}

	_, err = expression.Evaluate(nil)
	if err != nil {
		return errorz.ErrInvalidExpression
	}

	return nil
}
