package hserver

import (
	"encoding/json"
	"net/http"

	"github.com/thnxvlad/oplati/internal/domain"
)

type GetUsersResponse struct { 
	Users []domain.UserInfo `json:"users"`
 }

func (s *Server) getUsersInfoHandler(w http.ResponseWriter, r *http.Request) {
	users := s.oplatiService.GetUsersInfo(r.Context())
	response := GetUsersResponse{
		Users : users,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}