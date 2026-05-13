package hserver

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	hmiddlewares "github.com/thnxvlad/oplati/internal/server/hmiddlewares"
)
type GetUserResponse struct {
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

func (s *Server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(hmiddlewares.AccountIdContextKey{}).(uuid.UUID)

	ui, err := s.oplatiService.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := GetUserResponse{
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
