package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type AccountCreated struct {
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
		err = db.Append(&AccountCreated{
			Owner: "florimond",
		})
		assert.NoError(t, err)
	}

	changes, err := db.FindChanges(8, "account.created")
	assert.NoError(t, err)
	assert.Len(t, changes, 2)
	assert.Equal(t, []Event{
		&AccountCreated{Owner: "florimond"},
		&AccountCreated{Owner: "florimond"},
	}, changes)
}
