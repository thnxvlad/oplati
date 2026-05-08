package hserver

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type TransferRequest struct {
	SenderID uuid.UUID `json:"senderID"`
	RecipientID uuid.UUID `json:"recipientID"`
	Amount int `json:"amount"`
}


func (s *Server) transferHandler(w http.ResponseWriter, r *http.Request) {
	request := TransferRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.oplatiService.Transfer(r.Context(), request.SenderID, request.RecipientID, request.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
