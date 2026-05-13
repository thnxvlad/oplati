package hserver

import (
	"encoding/json"
	"net/http"
)

type NewUserRequest struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

func (s *Server) newUserHandler(w http.ResponseWriter, r *http.Request) {
	request := NewUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := s.oplatiService.CreateUser(r.Context(), request.Login, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
