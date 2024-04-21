package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/rs/zerolog/log"
)

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var loginRequest entity.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Error().Err(err).Msg("failed to parse request")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad json"})
		return
	}

	token, err := h.userUseCase.Login(ctx, loginRequest)
	switch {
	case errors.Is(err, errorz.ErrUserNotFound):
		log.Error().Err(err).Msg("user with this login and password not found")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "user not found"})
		return
	case err != nil:
		log.Error().Err(err).Msg("failed to login user")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}
