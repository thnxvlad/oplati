package oplati

import (
	"sync"

	"github.com/google/uuid"
	"github.com/thnxvlad/oplati/internal/domain"
)

type Storage struct {
	db map[uuid.UUID]domain.UserInfo
	sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		db: make(map[uuid.UUID]domain.UserInfo),
	}
}
