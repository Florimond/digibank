package transaction

// Use-case:
// As a current customer I would like to view my transactions so I can see my transaction history

// Manager represents a manager for the transactions
type Manager struct {
	db EventStore
}

// NewManager creates a new manager for transactions
func NewManager(db EventStore) (*Manager, error) {
	db.Register(eventTransaction, &Transaction{})
	return &Manager{
		db: db,
	}, nil
}

// Process processes commands
func (m *Manager) Process(command Command) error {
	switch command := command.(type) {
	case *DepositCommand:
		m.appendTx(&Transaction{
			AccountFrom: "ATM",
			AccountTo:   command.AccountTo,
			Amount:      command.Amount,
		})

	case *WithdrawCommand:
		// TODO: validate if AccountFrom has enough money for a transfer, otherwise error out
		m.appendTx(&Transaction{
			AccountFrom: command.AccountFrom,
			AccountTo:   "ATM",
			Amount:      command.Amount,
		})

	case *TransferCommand:
		// TODO: validate if AccountFrom has enough money for a transfer, otherwise error out
		m.appendTx(&Transaction{
			AccountFrom: command.AccountFrom,
			AccountTo:   command.AccountTo,
			Amount:      command.Amount,
		})
	}

	return nil
}

// appendTx adds a transaction to the database
func (m *Manager) appendTx(tx *Transaction) {
	m.db.Append(tx)

	// TODO: update snapshots
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
func (m *Manager) ViewBalance(account string) (float64, error) {
	//TODO
	return 0, nil
}
