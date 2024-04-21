package user

import (
	"encoding/json"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func (h Handler) UpdateOperation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	operations := map[string]string{
		"+": "addition",
		"-": "subtraction",
		"*": "multiplication",
		"/": "division",
	}

	var updateOperationRequest entity.UpdateOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&updateOperationRequest); err != nil {
		log.Error().Err(err).Msg("failed to parse request")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad json"})
		return
	}

	if err := ValidateUpdateOperation(updateOperationRequest); err != nil {
		log.Error().Err(err).Msg("failed to parse request")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad request"})
		return
	}

	updateOperationRequest.Operation = operations[updateOperationRequest.Operation]
	err := h.userUseCase.UpdateOperation(ctx, updateOperationRequest)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse request")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad request"})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ValidateUpdateOperation(request entity.UpdateOperationRequest) error {
	if !strings.Contains("+-*/", request.Operation) {
		return errorz.ErrInvalidOperation
	}
	if request.Time < 0 {
		return errorz.ErrInvalidOperationTime
	}
	return nil
}
