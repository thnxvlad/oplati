package hserver

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	hmiddlewares "github.com/thnxvlad/oplati/internal/server/hmiddlewares"
)

type TransferRequest struct {
	RecipientID uuid.UUID `json:"recipientID"`
	Amount      int       `json:"amount"`
}

func (s *PublicServer) transferHandler(w http.ResponseWriter, r *http.Request) {
	request := TransferRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(hmiddlewares.AccountIdContextKey{}).(uuid.UUID)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.oplatiService.Transfer(r.Context(), userID, request.RecipientID, request.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
