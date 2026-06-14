package oplati_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	oplatiStorage "github.com/thnxvlad/oplati/internal/storages/inmemory/oplati"
)

func TestCreateUser(t *testing.T) {
	t.Run("success create user", func(t *testing.T) {
		testCases := []struct {
			name string
			id   uuid.UUID
		}{
			{
				name: "random user 1",
				id:   uuid.New(),
			},
			{
				name: "random user 2",
				id:   uuid.New(),
			},
			{
				name: "random user 3",
				id:   uuid.New(),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				st := oplatiStorage.NewStorage()

				err := st.CreateUser(context.Background(), tc.id)
				if err != nil {
					t.Fatal(err)
				}

				user, err := st.GetUser(context.Background(), tc.id)
				if err != nil {
					t.Fatal("user does not exist")
				}

				if user.Balance != 0 {
					t.Fatalf("balance %d", user.Balance)
				}
			})
		}
	})

	t.Run("user id is already exist", func(t *testing.T) {
		testCases := []struct {
			name string
			id   uuid.UUID
		}{
			{
				name: "duplicate user 1",
				id:   uuid.New(),
			},
			{
				name: "duplicate user 2",
				id:   uuid.New(),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				st := oplatiStorage.NewStorage()

				err := st.CreateUser(context.Background(), tc.id)
				if err != nil {
					t.Fatal(err)
				}

				err = st.CreateUser(context.Background(), tc.id)
				if err == nil {
					t.Fatalf("error is nil")
				}

				if err.Error() != "id already exists" {
					t.Fatalf("unexpected error: %v", err)
				}
			})
		}
	})
}

func TestDeposit(t *testing.T) {
	testCases := []struct {
		name           string
		depositAmount  int
		expectedAmount int
	}{
		{
			name:           "deposit 30",
			depositAmount:  30,
			expectedAmount: 30,
		},
		{
			name:           "deposit 100",
			depositAmount:  100,
			expectedAmount: 100,
		},
		{
			name:           "deposit 999",
			depositAmount:  999,
			expectedAmount: 999,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			st := oplatiStorage.NewStorage()
			id := uuid.New()

			err := st.CreateUser(context.Background(), id)
			if err != nil {
				t.Fatal(err)
			}

			err = st.Deposit(context.Background(), id, tc.depositAmount)
			if err != nil {
				t.Fatal(err)
			}

			user, err := st.GetUser(context.Background(), id)
			if err != nil {
				t.Fatal(err)
			}

			if user.Balance != tc.expectedAmount {
				t.Fatalf(
					"user balance must be %d, got %d",
					tc.expectedAmount,
					user.Balance,
				)
			}
		})
	}
}
