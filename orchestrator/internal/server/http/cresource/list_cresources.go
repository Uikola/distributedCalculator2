package cresource

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (h Handler) ListCResources(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cResources, err := h.cResourceUseCase.ListCResources(ctx)
	if err != nil {
		log.Error().Err(err).Msg("can't get list of computing resources")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
		return
	}

	_ = json.NewEncoder(w).Encode(cResources)
}
