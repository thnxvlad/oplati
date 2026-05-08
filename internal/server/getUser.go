package hserver

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type GetUserRequest struct {
	ID uuid.UUID `json:"id"`
}

type GetUserResponse struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Balance int       `json:"balance"`
}

func (s *Server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	request := GetUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ui, err := s.oplatiService.GetUser(r.Context(), request.ID)
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
