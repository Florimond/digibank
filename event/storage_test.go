package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type AccountCreated struct {
	ID
	Owner string `json:"owner"`
}

func (e *AccountCreated) Name() string {
	return "account.created"
}

func TestStorage(t *testing.T) {
	db, err := Open("")
	assert.NoError(t, err)

	// Add few events into the storage
	for i := 0; i < 10; i++ {
		id, err := db.Append(&AccountCreated{
			Owner: "florimond",
		})
		assert.NoError(t, err)
		assert.NotEqual(t, 0, id)
	}

	changes, err := db.FindChanges(8, "account.created")
	assert.NoError(t, err)
	assert.Len(t, changes, 2)
	assert.Equal(t, []Event{
		&AccountCreated{Owner: "florimond", ID: ID{9}},
		&AccountCreated{Owner: "florimond", ID: ID{10}},
	}, changes)
}
