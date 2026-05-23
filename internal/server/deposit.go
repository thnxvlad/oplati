package hserver

import (
	"encoding/json"
	"net/http"

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

	userID, err := hmiddlewares.GetAccountIdFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = s.oplatiService.Deposit(r.Context(), userID, request.Dep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
