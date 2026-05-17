package hserver

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	hmiddlewares "github.com/thnxvlad/oplati/internal/server/hmiddlewares"
)

type DepositRequest struct {
	Dep int `json:"deposit"`
}

func (s *PublicServer) depositHandler(w http.ResponseWriter, r *http.Request) {
	request := DepositRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, _ := r.Context().Value(hmiddlewares.AccountIdContextKey{}).(uuid.UUID)

	err = s.oplatiService.Deposit(r.Context(), userID, request.Dep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
