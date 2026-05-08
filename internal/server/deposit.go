package hserver

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type DepositRequest struct {
	ID uuid.UUID `json:"id"`
	Dep int `json:"deposit"`
}

type DepositResponse struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Balance int       `json:"balance"`
}

func (s *Server) depositHandler(w http.ResponseWriter, r *http.Request) {
	request := DepositRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ui, err := s.oplatiService.Deposit(r.Context(), request.ID, request.Dep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := DepositResponse{
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
