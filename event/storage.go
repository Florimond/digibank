package event

import (
	"fmt"
	"reflect"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Storage abstracts the event sourcing database
type Storage struct {
	db    *gorm.DB
	types map[string]reflect.Type
}

// Open opens the database
func Open(name string) (*Storage, error) {

	// If no name was specified, open in memory (for testing)
	if name == "" {
		name = "file::memory:?cache=shared"
	}

	// Open the SQLite database
	db, err := gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&record{})
	return &Storage{
		db:    db,
		types: make(map[string]reflect.Type),
	}, nil
}

// Append appends an event into the store
func (s *Storage) Append(event Event) (uint, error) {
	s.Register(event.Name(), event)
	newRec := newRecord(event)
	err := s.db.Create(newRec).Error
	event.SetEventID(newRec.ID)
	return newRec.ID, err
}

// FindChanges finds all of the changes after a certain key
func (s *Storage) FindChanges(after uint, names ...string) ([]Event, error) {
	records := []record{}
	if tx := s.db.Debug().
		Order("id").
		Where("name IN ? AND id > ?", names, after).
		Find(&records); tx.Error != nil {
		return nil, tx.Error
	}

	return s.makeEvents(records)
}

// Register registers a type of event into the store so we can create it
// while querying
func (s *Storage) Register(name string, event Event) {
	if _, ok := s.types[name]; !ok {
		s.types[name] = reflect.TypeOf(event).Elem()
	}
}

// makeEvent creates an instance of an event from a record
func (s *Storage) makeEvent(r *record) (Event, error) {
	if typ, ok := s.types[r.Name]; ok {
		return reflect.New(typ).Interface().(Event), nil
	}

	return nil, fmt.Errorf("event: unknown type %s", r.Name)
}

// makeEvents convert records to events
func (s *Storage) makeEvents(records []record) ([]Event, error) {
	result := make([]Event, 0, len(records))
	for _, r := range records {
		event, err := s.makeEvent(&r)
		if err != nil {
			return nil, err
		}

		event.SetEventID(r.ID)
		if err := r.Unmarshal(event); err != nil {
			return nil, err
		}

		result = append(result, event)
	}
	return result, nil
}
