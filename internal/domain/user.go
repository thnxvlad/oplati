package domain

import "github.com/google/uuid"

type UserInfo struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Balance int       `json:"balance"`
	Login string
	Password string
}

const JwtSecret = "fdmslfnsdf"
