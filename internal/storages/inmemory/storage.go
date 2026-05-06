package inmemory

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

func (s *Storage) CreateUser(ctx context.Context, ui domain.UserInfo) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.db[ui.Id]; ok {
		return errors.New("id already exists")
	}

	s.db[ui.Id] = ui

	return nil
}

func (s *Storage) Deposit(ctx context.Context, userId uuid.UUID, amount int) (domain.UserInfo, error) {
	s.Lock()
	defer s.Unlock()

	ui, ok := s.db[userId]
	if !ok {
		return domain.UserInfo{}, errors.New("user id does not exist")
	}

	ui.Balance += amount

	s.db[userId] = ui

	return ui, nil
}

func (s *Storage) GetUser(ctx context.Context, userId uuid.UUID) (domain.UserInfo, error) {
	s.RLock()
	defer s.RUnlock()

	ui, ok := s.db[userId]
	if !ok {
		return domain.UserInfo{}, errors.New("user id does not exist")
	}

	return ui, nil
}

func (s *Storage) GetUsersInfo(ctx context.Context) ([]domain.UserInfo, error) {
	s.RLock()
	defer s.RUnlock()
	var users []domain.UserInfo
	for _, value := range s.db{
			users = append(users, value)
	} 
	//пока не придумал когда ошибку выводить
	return users, nil
}