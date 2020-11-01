package account

import (
	"testing"

	"github.com/florhusq/digibank/event"
)

func TestCommands(t *testing.T) {
	db, err := event.Open("")
	if err != nil {
		t.Fatal(err)
	}

	manager, err := NewManager(db)
	if err != nil {
		t.Fatal(err)
	}

	err = manager.createAccount("florimond")
	if err != nil {
		t.Fatal(err)
	}

	//manager.findAccount("")

}
