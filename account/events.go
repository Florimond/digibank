package account

import (
	"github.com/florhusq/digibank/event"
)

const (
	eventTransaction = "transaction"
	eventOpenAccount = "openAccount"
)

// EventStore represents an event source (dependency inversion principle)
type EventStore interface {
	Append(event event.Event) (uint, error)
	FindChanges(after uint, names ...string) ([]event.Event, error)
	Register(name string, event event.Event)
}

// Transaction represents a transaction event
type Transaction struct {
	event.ID
	AccountFrom string  `json:"from"`   // The account which sends the money
	AccountTo   string  `json:"to"`     // The account which receives the money
	Amount      float64 `json:"amount"` // The amount
}

// Name returns the event name
func (t *Transaction) Name() string {
	return eventTransaction
}

// OpenAccount represents the opening of an account
type OpenAccount struct {
	event.ID
	AccountID string `json:"account"`  // The ID of the new account
	Customer  string `json:"customer"` // The customer owning the new account
}

// Name returns the event name
func (oa *OpenAccount) Name() string {
	return eventOpenAccount
}
