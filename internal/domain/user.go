package domain

import "github.com/google/uuid"

type UserInfo struct {
	Id      uuid.UUID `json:"id"`
	Balance int       `json:"balance"`
}
