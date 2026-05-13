package auth

import (
	"context"

	"github.com/google/uuid"
)

type AuthDB interface {
	// ToDo: add methods
}

type OplatiService interface {
	// для того чтобы создать пользователя в базе данных oplati, позде перепишем под использование брокера сообщений
	CreateUser(ctx context.Context, userID uuid.UUID) error
}

type Service struct {
	OplatiService OplatiService
	db            AuthDB
}

func New(db AuthDB, oplatiService OplatiService) *Service {
	return &Service{
		db:            db,
		OplatiService: oplatiService,
	}
}
