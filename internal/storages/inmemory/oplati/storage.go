package oplati

import (
	"context"
	"errors"
	"fmt"
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
		return nil, errors.New("context cancelled")
	}
	s.RLock()
	defer s.RUnlock()
	var users []domain.UserInfo
	for _, value := range s.db {
		users = append(users, value)
	}

	return users, nil
}
type balanceUpdate struct {
	id uuid.UUID
  	amount int
}

func (s *Storage) updateBalance(userIDs []balanceUpdate) error {
	users := make([]domain.UserInfo, len(userIDs))
	var ok bool 
	s.Lock()
	defer s.Unlock()
	for i, value := range userIDs {
		users[i], ok = s.db[value.id]
		if !ok {
			return fmt.Errorf("user %s does not exist", value.id.String())
		}

		users[i].Balance += value.amount
		if users[i].Balance < 0 {
			return errors.New("not enough funds")
		}
	}
	for _, value := range users{
		s.db[value.Id] = value
	}
	

	return nil
}

func (s *Storage) Transfer(ctx context.Context, senderID uuid.UUID, recipientID uuid.UUID, amount int) error {
	if ctx.Err() != nil {
		return errors.New("context cancelled")
	} 
	
	si := balanceUpdate{
		id : senderID,
		amount: amount,
	}

	ri := balanceUpdate{
		id : recipientID,
		amount: -amount,
	}

	s.updateBalance([]balanceUpdate{si,ri})

	return nil
}

func (s *Storage) Deposit(ctx context.Context, userID uuid.UUID, amount int) error {
	if ctx.Err() != nil {
		return errors.New("context cancelled")
	} 
	
	ui := balanceUpdate{
		id : userID,
		amount: amount,
	}

	s.updateBalance([]balanceUpdate{ui})

	return nil
}

func (s *Storage) Withdraw(ctx context.Context, userID uuid.UUID, amount int) error {
	if ctx.Err() != nil {
		return errors.New("context cancelled")
	} 
	
	ui := balanceUpdate{
		id : userID,
		amount: -amount,
	}

	s.updateBalance([]balanceUpdate{ui})

	return nil
}

func (s *Storage) GetUser(ctx context.Context, userId uuid.UUID) (domain.UserInfo, error) {
	if ctx.Err() != nil {
		return domain.UserInfo{}, errors.New("context cancelled")
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
		return errors.New("context cancelled")
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

