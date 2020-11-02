package account

import (
	"testing"

	"github.com/florhusq/digibank/event"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) *Manager {
	db, err := event.Open("")
	if err != nil {
		t.Fatal(err)
	}

	manager, err := NewManager(db)
	if err != nil {
		t.Fatal(err)
	}

	return manager
}

func Test_createAccount(t *testing.T) {
	manager := setup(t)

	accID, err := manager.createAccount("florimond")
	if err != nil {
		t.Fatal(err)
	}

	acc, err := manager.findAccount(accID)
	assert.Nil(t, err)
	assert.Equal(t, "florimond", acc.Customer)
}

func Test_transfers(t *testing.T) {
	manager := setup(t)

	// Create an account
	openAccount1 := &OpenAccountCommand{
		Customer: "florimond",
	}
	accFlorimondID, err := manager.Process(openAccount1)
	if err != nil {
		t.Fatal(err)
	}

	// Deposit to this account
	depositToFlo := &DepositCommand{
		AccountTo: accFlorimondID,
		Amount:    50.0,
	}
	manager.Process(depositToFlo)

	// Check the balance after this deposit
	balance, err := manager.ViewBalance(accFlorimondID)
	assert.Nil(t, err)
	assert.Equal(t, 50.0, balance)

	// Withdraw from thia account
	withdrawFlo := &WithdrawCommand{
		AccountFrom: accFlorimondID,
		Amount:      25.0,
	}
	manager.Process(withdrawFlo)

	// Check the balance after this withdrawal
	balance, err = manager.ViewBalance(accFlorimondID)
	assert.Nil(t, err)
	assert.Equal(t, 25.0, balance)

	// Create another account for a transfer
	openAccount := &OpenAccountCommand{
		Customer: "emilie",
	}
	accEmilieID, err := manager.Process(openAccount)
	if err != nil {
		t.Fatal(err)
	}

	// Transfer from one account to another.
	transferToEmi := &TransferCommand{
		AccountFrom: accFlorimondID,
		AccountTo:   accEmilieID,
		Amount:      25,
	}
	manager.Process(transferToEmi)

	// Check the balance on both accounts.
	balance, err = manager.ViewBalance(accFlorimondID)
	assert.Nil(t, err)
	assert.Equal(t, 0.0, balance)

	balance, err = manager.ViewBalance(accEmilieID)
	assert.Nil(t, err)
	assert.Equal(t, 25.0, balance)
}
