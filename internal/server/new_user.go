package hserver

import (
	"net/http"

	"github.com/google/uuid"
)

type NewUserRequest struct {
	Name string `json:"name"`
}

type NewUserResponse struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Balance string    `json:"balance"`
}

func NewUserFunc(w http.ResponseWriter, r *http.Request) {

}
