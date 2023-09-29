package store

import (
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------------------------------- //
// Value store interface and methods

//
// ValueStore - Getter/Setter for string map of interfaces{}
// ---------------------------------------------------------------------------------------------------- //
type ValueStore interface {
	// NewValueStore - Set up the ValueStore map.
	NewValueStore()

	// ValueExists - Check if a key exists.
	ValueExists(key string) bool

	// ValueNotExists - Inverse of Exists()
	ValueNotExists(key string) bool

	// GetValue - Get a key's value.
	GetValue(key string) any

	// SetValue - Set a key value pair.
	SetValue(key string, value any)

	// CountValues - Return the number of entries.
	CountValues() int

	// String - Stringer method.
	String() string
}

// NewValueStore - Create a ValueStore interface structure instance.
//goland:noinspection GoUnusedExportedFunction
func NewValueStore() ValueStore {
	return &ValueStruct{
		Values: make(map[string]any),
	}
}

//
// ValueStruct
// ---------------------------------------------------------------------------------------------------- //
type ValueStruct struct {
	Values map[string]any `json:"values"`
}

// NewValueStruct - Create a ValueStore interface structure instance.
func NewValueStruct() ValueStruct {
	return ValueStruct{
		Values: make(map[string]any),
	}
}

// NewValueStore - Create a ValueStore interface structure instance.
func (p *ValueStruct) NewValueStore() {
	p.Values = make(map[string]any)
}

// ValueExists - Check if a key exists.
func (p *ValueStruct) ValueExists(key string) bool {
	key = strings.TrimSpace(key)
	if _, ok := p.Values[key]; ok {
		return true
	}
	return false
}

// ValueNotExists - Inverse of ValueExists()
func (p *ValueStruct) ValueNotExists(key string) bool {
	key = strings.TrimSpace(key)
	if _, ok := p.Values[key]; ok {
		return false
	}
	return true
}

// GetValue - Get a key's value.
func (p *ValueStruct) GetValue(key string) any {
	key = strings.TrimSpace(key)
	if value, ok := p.Values[key]; ok {
		return value
	}
	return new(any)
}

// SetValue - Set a key value pair.
func (p *ValueStruct) SetValue(key string, value any) {
	key = strings.TrimSpace(key)
	p.Values[key] = value
}

// CountValues - Return the number of entries.
func (p *ValueStruct) CountValues() int {
	return len(p.Values)
}

// String - Stringer interface.
func (p ValueStruct) String() string {
	var ret string
	for key, value := range p.Values {
		ret += fmt.Sprintf("ValueStruct[%s] => %v\n",
			key, value)
	}
	return ret
}
