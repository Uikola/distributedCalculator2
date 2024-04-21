package user

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func (h Handler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Error().Msg("auth header is empty")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"reason": "auth header is empty"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Error().Msg("invalid auth header")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"reason": "invalid auth header"})
			return
		}

		userID, err := h.userUseCase.ParseToken(r.Context(), parts[1])
		if err != nil {
			log.Error().Err(err).Msg("invalid auth token")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"reason": "invalid auth token"})
			return
		}

		ctx := context.WithValue(context.WithValue(r.Context(), "userID", userID), "token", parts[1])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
