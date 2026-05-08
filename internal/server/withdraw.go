package hserver

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type WithdrawRequest struct {
	ID     uuid.UUID `json:"id"`
	Amount int       `json:"deposit"`
}

type WithdrawResponse struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Balance int       `json:"balance"`
}

func (s *Server) withdrawHandler(w http.ResponseWriter, r *http.Request) {
	request := WithdrawRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if request.Amount < 0 {
		http.Error(w, "amount must be positive", http.StatusBadRequest)
		return
	}
	ui, err := s.oplatiService.Withdraw(r.Context(), request.ID, request.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := WithdrawResponse{
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
