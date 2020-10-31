package transaction

import (
	"github.com/florhusq/digibank/event"
)

const (
	eventTransaction = "transaction"
)

// EventStore represents an event source (dependency inversion principle)
type EventStore interface {
	Append(event event.Event) error
	FindChanges(after uint, names ...string) ([]event.Event, error)
	Register(name string, event event.Event)
}

// Transaction represents a transaction event
type Transaction struct {
	AccountFrom string  `json:"from"`   // The account which sends the money
	AccountTo   string  `json:"to"`     // The account which receives the money
	Amount      float64 `json:"amount"` // The amount
}

// Name returns the event name
func (t *Transaction) Name() string {
	return eventTransaction
}

// eventsToTransactions converts events to transaction
func eventsToTransactions(events []event.Event) []Transaction {
	result := make([]Transaction, 0, len(events))
	for _, v := range events {
		result = append(result, *v.(*Transaction))
	}
	return result
}
