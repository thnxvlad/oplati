package hserver

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type SignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

func (s *PrivateServer) signInHandler(w http.ResponseWriter, r *http.Request) {
	req := SignInRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "login and password are required", http.StatusBadRequest)
		return
	}

	token, err := s.authService.SignIn(req.Login, req.Password)
	if err != nil {
		log.Error().Err(err).Str("login", req.Login).Msg("failed to sign in in auth service")
		http.Error(w, "login failed", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(SignInResponse{Token: token}); err != nil {
		log.Error().Err(err).Msg("failed to encode response")
	}
}
