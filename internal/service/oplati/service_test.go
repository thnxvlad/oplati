package oplati

import (
	"context"
	"testing"

	"github.com/google/uuid"
	oplatiStorage "github.com/thnxvlad/oplati/internal/storages/inmemory/oplati"
)

func TestCreateUser(t *testing.T) {
	t.Run("success create user", func(t *testing.T) {
		st := oplatiStorage.NewStorage()
		id := uuid.New()

		err := st.CreateUser(context.Background(), id)
		if err != nil {
			t.Fatal(err)
		}

		user, err := st.GetUser(context.Background(), id)
		if err != nil {
			t.Fatal("user does not exist")
		}

		if user.Balance != 0 {
			t.Fatalf("balance %d", user.Balance)
		}
	})

	t.Run("user id is already exist", func(t *testing.T) {
		st := oplatiStorage.NewStorage()
		id := uuid.New()

		err := st.CreateUser(context.Background(), id)
		if err != nil {
			t.Fatal(err)
		}

		err = st.CreateUser(context.Background(), id)
		if err == nil {
			t.Fatalf("error is nil")
		}

		if err.Error() != "id already exists" {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestDeposit(t *testing.T) {
	t.Run("successful deposit", func(t *testing.T) {
		st := oplatiStorage.NewStorage()
		id := uuid.New()

		err := st.CreateUser(context.Background(), id)
		if err != nil {
			t.Fatal(err)
		}

		err = st.Deposit(context.Background(), id, 30)

		if err != nil {
			t.Fatal(err)
		}

		user, err := st.GetUser(context.Background(), id)
		

		if user.Balance != 30 {
			t.Fatalf("user balance must be 30, got %d", user.Balance)
		}

	})

}
