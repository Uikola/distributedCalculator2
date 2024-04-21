package expression

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func (h Handler) GetExpression(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	expressionID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("invalid task id")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "invalid task id"})
		return
	}

	userID := ctx.Value("userID").(uint)

	expression, err := h.expressionUseCase.GetExpression(ctx, userID, uint(expressionID))
	switch {
	case errors.Is(errorz.ErrExpressionNotFound, err):
		log.Error().Err(err).Msg("expression not found")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "expression not found"})
		return
	case errors.Is(errorz.ErrAccessForbidden, err):
		log.Error().Err(err).Msg("access for this expression forbidden")
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "access for this expression forbidden"})
		return
	case err != nil:
		log.Error().Msg(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(expression)
}
