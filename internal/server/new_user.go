package hserver

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type NewUserRequest struct {
	Name string `json:"name"`
}

type NewUserResponse struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Balance int       `json:"balance"`
}

func (s *Server) newUserHandler(w http.ResponseWriter, r *http.Request) {
	request := NewUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ui, err := s.oplatiService.CreateUser(r.Context(), request.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := NewUserResponse{
		Id:      ui.Id,
		Name:    ui.Name,
		Balance: ui.Balance,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
