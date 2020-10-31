package event

import (
	"encoding/json"
	"gorm.io/gorm"
)

// Event represents an event
type Event interface {
	Name() string
}

// Record represents an event stored in the database
type record struct {
	gorm.Model
	Name string `gorm:"index"`      // Name/type of the event
	Data []byte `gorm:"size:65536"` // JSON payload
}

// Unmarshal unmarshals the value into the destination
func (e *record) Unmarshal(dst Event) error {
	return json.Unmarshal(e.Data, dst)
}

// newRecord creates a new event
func newRecord(event Event) *record {
	b, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}

	return &record{
		Name: event.Name(),
		Data: b,
	}
}
