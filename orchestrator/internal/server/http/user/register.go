package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"unicode/utf8"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/entity"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/errorz"
	"github.com/rs/zerolog/log"
)

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var registerRequest entity.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		log.Error().Err(err).Msg("failed to parse request")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad json"})
		return
	}

	if err := Validate(registerRequest); err != nil {
		log.Error().Err(err).Msg("failed to validate request")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "bad request"})
		return
	}

	err := h.userUseCase.Create(ctx, registerRequest)
	switch {
	case errors.Is(err, errorz.ErrUserAlreadyExists):
		log.Error().Err(err).Msg("user with this login already exists")
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "user with this login already exists"})
		return
	case err != nil:
		log.Error().Msg(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"reason": "internal error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"response": "user created successfully"})
}

func Validate(user entity.RegisterRequest) error {
	if err := ValidateLogin(user.Login); err != nil {
		return err
	}
	if err := ValidatePassword(user.Password); err != nil {
		return err
	}
	return nil
}

func ValidateLogin(login string) error {
	if utf8.RuneCountInString(login) < 4 || utf8.RuneCountInString(login) > 20 {
		return errorz.ErrInvalidLogin
	}
	return nil
}

func ValidatePassword(password string) error {
	var upper, lower, digit bool
	for _, el := range password {
		if el > 96 && el < 123 {
			lower = true
		}
		if el > 64 && el < 91 {
			upper = true
		}
		if el > 47 && el < 58 {

			digit = true
		}
	}

	if utf8.RuneCountInString(password) < 6 || utf8.RuneCountInString(password) > 100 || !upper || !lower || !digit {
		return errorz.ErrInvalidPassword
	}
	return nil
}
