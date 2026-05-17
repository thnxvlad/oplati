package hserver

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	hmiddlewares "github.com/thnxvlad/oplati/internal/server/hmiddlewares"
)

type WithdrawRequest struct {
	Amount int `json:"deposit"`
}

func (s *PublicServer) withdrawHandler(w http.ResponseWriter, r *http.Request) {
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

	userID, _ := r.Context().Value(hmiddlewares.AccountIdContextKey{}).(uuid.UUID)

	err = s.oplatiService.Withdraw(r.Context(), userID, request.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
