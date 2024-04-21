package expression

import (
	"encoding/json"
	"errors"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h Handler) ListExpressions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	expressions, err := h.expressionUseCase.ListExpressions(ctx)
	switch {
	case errors.Is(errorz.ErrNoExpressions, err):
		log.Error().Err(err).Msg("user doesn't have any expressions")
		w.WriteHeader(http.StatusNoContent)
		_ = json.NewEncoder(w).Encode(map[string]string{"response": "you don't have any expressions"})
		return
	case err != nil:
		log.Error().Msg(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(expressions)
}
