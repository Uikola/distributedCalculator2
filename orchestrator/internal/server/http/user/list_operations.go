package user

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (h Handler) ListOperations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	operations, err := h.userUseCase.ListOperations(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get the list of user's operatons")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(operations)
}
