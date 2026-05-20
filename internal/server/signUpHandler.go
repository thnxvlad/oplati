package hserver

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

func (s *PrivateServer) registerHandler(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())

	req := RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to decode request")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		logger.Warn().Msg("empty login or password")
		http.Error(w, "login and password are required", http.StatusBadRequest)
		return
	}

	user, err := s.oplatiService.CreateUser(r.Context(), req.Name)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create user in storage")
		http.Error(w, "could not create user", http.StatusInternalServerError)
		return
	}

	err = s.authService.SignUp(req.Login, req.Password, user.Id.String())
	if err != nil {
		logger.Error().Err(err).Str("login", req.Login).Msg("failed to sign up in auth service")
		http.Error(w, "registration failed", http.StatusBadRequest)
		return
	}

	token, err := s.authService.SignIn(req.Login, req.Password)
	if err != nil {
		logger.Error().Err(err).Msg("failed to login after registration")
		http.Error(w, "auth failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(RegisterResponse{Token: token}); err != nil {
		logger.Error().Err(err).Msg("failed to encode response")
	}
}
