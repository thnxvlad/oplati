package oplati

import (
	"context"
	"errors"
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

func (s *Storage) GetUsersInfo(ctx context.Context) ([]domain.UserInfo, error) {
	if ctx.Err() != nil {
		return nil, errors.New("request canceled")
	}
	s.RLock()
	defer s.RUnlock()
	var users []domain.UserInfo
	for _, value := range s.db {
		users = append(users, value)
	}

	return users, nil
}

func (s *Storage) UpdateBalance(ctx context.Context, userId uuid.UUID, amount int) error {
	if ctx.Err() != nil {
		return errors.New("context cancelled")
	}
	s.Lock()
	defer s.Unlock()

	ui, ok := s.db[userId]
	if !ok {
		return errors.New("user id does not exist")
	}

	ui.Balance += amount
	if ui.Balance < 0 {
		return errors.New("not enough money")
	}
	s.db[userId] = ui

	return nil
}

// func (s *Storage) updateBalance(ctx context.Context, userId uuid.UUID, amount int) error {
// 	ui, ok := s.db[userId]
// 	if !ok {
// 		return errors.New("user id does not exist")
// 	}
// 	ui.Balance += amount
// 	if ui.Balance < 0 {
// 		return errors.New("not enough money")
// 	}
// 	s.db[userId] = ui

// 	return nil
// }

func (s *Storage) Transfer(ctx context.Context, senderID uuid.UUID, recipientID uuid.UUID, amount int) error {
	if ctx.Err() != nil {
		return errors.New("context cancelled")
	} 
	s.Lock()
	defer s.Unlock()

	si, ok := s.db[senderID]
	if !ok {
		return errors.New("sender id does not exist")
	}
	ri, ok := s.db[recipientID]
	if !ok {
		return errors.New("recipient id does not exist")
	}

	si.Balance -= amount
	if si.Balance < 0 {
		return errors.New("not enough funds")
	}
	s.db[senderID] = si

	ri.Balance += amount
	s.db[recipientID] = ri

	return nil
}

func (s *Storage) GetUser(ctx context.Context, userId uuid.UUID) (domain.UserInfo, error) {
	if ctx.Err() != nil {
		return domain.UserInfo{}, errors.New("request terminated")
	} 
	s.RLock()
	defer s.RUnlock()

	ui, ok := s.db[userId]
	if !ok {
		return domain.UserInfo{}, errors.New("user id does not exist")
	}

	return ui, nil
}

func (s *Storage) CreateUser(ctx context.Context, userId uuid.UUID) error {
	if ctx.Err() != nil {
		return errors.New("request terminated")
	} 
	s.Lock()
	defer s.Unlock()

	if _, ok := s.db[userId]; ok {
		return errors.New("id already exists")
	}

	ui := domain.UserInfo{
		Id:       userId,
		Balance:  0,
	}

	s.db[ui.Id] = ui
	return nil
}

