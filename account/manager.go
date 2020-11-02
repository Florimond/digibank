package account

import (
	"errors"
	"sync"

	"github.com/florhusq/digibank/event"
	"github.com/google/uuid"
)

// Use-case:
// As a current customer I would like to view my transactions so I can see my transaction history

var errNoAccount = errors.New("account not found")
var errInsufficientFunds = errors.New("insufficient funds")

// Manager represents a manager for the transactions
type Manager struct {
	lock     sync.Mutex
	db       EventStore
	accounts map[string]*Account
}

// NewManager creates a new manager for transactions
func NewManager(db EventStore) (*Manager, error) {
	db.Register(eventTransaction, &Transaction{})
	db.Register(eventOpenAccount, &OpenAccount{})
	m := &Manager{
		db:       db,
		accounts: make(map[string]*Account, 0),
	}
	// Replay all the changes to rebuild the database
	m.ApplyChanges()
	return m, nil
}

// Process processes commands
func (m *Manager) Process(command Command) (string, error) {
	switch command := command.(type) {
	case *DepositCommand:
		return m.deposit(command)
	case *WithdrawCommand:
		return m.withdraw(command)
	case *TransferCommand:
		return m.transfer(command)
	case *OpenAccountCommand:
		return m.createAccount(command.Customer)
	}

	return "", nil
}

// transfer is the command that transfers money from an account to another
func (m *Manager) transfer(command *TransferCommand) (string, error) {
	if _, err := m.findAccount(command.AccountTo); err != nil {
		return "", errNoAccount
	}
	accFrom, err := m.findAccount(command.AccountFrom)
	if err != nil {
		return "", errNoAccount
	}
	if accFrom.Amount < command.Amount {
		return "", errInsufficientFunds
	}

	m.appendTx(&Transaction{
		AccountFrom: command.AccountFrom,
		AccountTo:   command.AccountTo,
		Amount:      command.Amount,
	})

	return "", nil
}

// withdraw is the command that withdraws money from the account
func (m *Manager) withdraw(command *WithdrawCommand) (string, error) {
	acc, err := m.findAccount(command.AccountFrom)
	if err != nil {
		return "", errNoAccount
	}
	if acc.Amount < command.Amount {
		return "", errInsufficientFunds
	}

	m.appendTx(&Transaction{
		AccountFrom: command.AccountFrom,
		AccountTo:   "ATM",
		Amount:      command.Amount,
	})

	return "", nil
}

// deposit is the command that deposits money into an account
func (m *Manager) deposit(command *DepositCommand) (string, error) {
	if _, err := m.findAccount(command.AccountTo); err != nil {
		return "", errNoAccount
	}

	m.appendTx(&Transaction{
		AccountFrom: "ATM",
		AccountTo:   command.AccountTo,
		Amount:      command.Amount,
	})

	return "", nil
}

// createAccount is the command that creates an account.
func (m *Manager) createAccount(customer string) (string, error) {
	event := &OpenAccount{
		AccountID: uuid.New().String(),
		Customer:  customer,
	}

	_, err := m.db.Append(event)
	if err != nil {
		return "", err
	}

	m.Apply(event)

	return event.AccountID, nil
}

// findAccount finds an account baed on its ID.
func (m *Manager) findAccount(ID string) (*Account, error) {
	acc, ok := m.accounts[ID]
	if !ok {
		return nil, errNoAccount
	}
	return acc, nil
}

// appendTx adds a transaction to the database
func (m *Manager) appendTx(tx *Transaction) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	_, err := m.db.Append(tx)
	if err != nil {
		return err
	}

	m.Apply(tx)

	return nil
}

// ViewTransactions shows all of the transactions for a user
func (m *Manager) ViewTransactions(account string) ([]Transaction, error) {
	events, err := m.db.FindChanges(0, eventTransaction)
	if err != nil {
		return nil, err
	}

	// Filter by account name
	result := []Transaction{}
	for _, v := range events {
		tx := *v.(*Transaction)
		if tx.AccountFrom == account || tx.AccountTo == account {
			result = append(result, *v.(*Transaction))
		}
	}
	return result, nil
}

// ViewBalance shows the balance of the account
func (m *Manager) ViewBalance(accountID string) (float64, error) {
	acc, err := m.findAccount(accountID)
	if err != nil {
		return 0, err
	}
	return acc.Amount, nil
}

// Apply applies the event received
func (m *Manager) Apply(e event.Event) {
	switch e := e.(type) {
	case *OpenAccount:
		m.accounts[e.AccountID] = &Account{
			ID:       e.AccountID,
			Customer: e.Customer,
			Amount:   0,
			Version:  e.EventID,
		}
	case *Transaction:
		// Error is ignored because the account is supposed to be existing at this stage
		if e.AccountFrom != "ATM" {
			accFrom, _ := m.findAccount(e.AccountFrom)
			accFrom.Amount -= e.Amount
			accFrom.Version = e.EventID
		}

		if e.AccountTo != "ATM" {
			accTo, _ := m.findAccount(e.AccountTo)
			accTo.Amount += e.Amount
			accTo.Version = e.EventID
		}
	}
}

// ApplyChanges replays all the changes since the beginning of times.
func (m *Manager) ApplyChanges() {
	events, err := m.db.FindChanges(0, eventOpenAccount, eventTransaction)
	if err != nil {
		panic(err)
	}
	for _, e := range events {
		m.Apply(e)
	}
}
