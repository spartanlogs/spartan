package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrKeyNotFound is returned during an operation when a given key doesn't exist
	ErrKeyNotFound = errors.New("key not found")
)

// An InterfaceMap is a wrapper object around map[string]interface{}.
// It's purpose is to more easily manipulate interface{} values in particular
// with multi-level data structures.
type InterfaceMap map[string]interface{}

// NewInterfaceMap creates and returns an InterfaceMap.
func NewInterfaceMap() InterfaceMap {
	return make(InterfaceMap)
}

// NewMap creates an InterfaceMap using an existing map.
func NewMap(data map[string]interface{}) InterfaceMap {
	return InterfaceMap(data)
}

// Set the key to val.
func (m InterfaceMap) Set(key string, val interface{}) {
	walkMap(key, m, mapOperation{putOperator{value: val}, true})
}

// Get the value of key.
func (m InterfaceMap) Get(key string) interface{} {
	v, _ := walkMap(key, m, mapOperation{getOperator{}, false})
	return v
}

// GetOK gets the value of key and if it exists.
func (m InterfaceMap) GetOK(key string) (interface{}, bool) {
	v, err := walkMap(key, m, mapOperation{getOperator{}, false})
	exists := (err != ErrKeyNotFound)
	return v, exists
}

// HasKey returns if the key exists in the map.
func (m InterfaceMap) HasKey(key string) bool {
	exists, _ := walkMap(key, m, mapOperation{hasKeyOperator{}, false})
	return exists.(bool)
}

// Delete key from the map.
func (m InterfaceMap) Delete(key string) error {
	_, err := walkMap(key, m, mapOperation{deleteOperator{}, false})
	return err
}

// Len returns the count of items in the map.
func (m InterfaceMap) Len() int {
	return len(m)
}

// Copy returns a new InterfaceMap with copies of the
// keys and values of m.
func (m InterfaceMap) Copy() InterfaceMap {
	n := NewInterfaceMap()
	for k, v := range m {
		innerMap, err := toInterfaceMap(v)
		if err == nil {
			n[k] = innerMap.Copy()
		} else {
			n[k] = v
		}
	}
	return n
}

// MarshalJSON marshals the underlying map instead of the struct.
func (m InterfaceMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalJSON unmarshals the underlying map instead of the struct.
func (m InterfaceMap) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, m)
}

func toInterfaceMap(data interface{}) (InterfaceMap, error) {
	switch m := data.(type) {
	case InterfaceMap:
		return m, nil
	case map[string]interface{}:
		return InterfaceMap(m), nil
	default:
		return nil, fmt.Errorf("Expected map, got %T", data)
	}
}

func walkMap(key string, data InterfaceMap, op mapOperation) (interface{}, error) {
	var err error
	keyPath := strings.Split(key, ".")

	m := data
	for i, k := range keyPath[:len(keyPath)-1] {
		v, exists := m[k]
		if !exists {
			if op.CreateMissingKeys {
				nm := InterfaceMap{}
				m[k] = nm
				m = nm
				continue
			}
			return nil, fmt.Errorf("key not found: %s", k)
		}

		m, err = toInterfaceMap(v)
		if err != nil {
			return nil, fmt.Errorf("invalid key %s", strings.Join(keyPath[:i+1], "."))
		}
	}

	v, err := op.Do(keyPath[len(keyPath)-1], m)
	if err != nil {
		return nil, fmt.Errorf("key=%s %s", key, err)
	}
	return v, nil
}

type mapOperation struct {
	mapOperator
	CreateMissingKeys bool
}

type mapOperator interface {
	Do(key string, data InterfaceMap) (interface{}, error)
}

type deleteOperator struct{}

func (op deleteOperator) Do(key string, data InterfaceMap) (interface{}, error) {
	v, exists := data[key]
	if !exists {
		return nil, ErrKeyNotFound
	}
	delete(data, key)
	return v, nil
}

type getOperator struct{}

func (op getOperator) Do(key string, data InterfaceMap) (interface{}, error) {
	v, exists := data[key]
	if !exists {
		return nil, ErrKeyNotFound
	}
	return v, nil
}

type hasKeyOperator struct{}

func (op hasKeyOperator) Do(key string, data InterfaceMap) (interface{}, error) {
	_, exists := data[key]
	return exists, nil
}

type putOperator struct {
	value interface{}
}

func (op putOperator) Do(key string, data InterfaceMap) (interface{}, error) {
	cv := data[key]
	data[key] = op.value
	return cv, nil
}
