package utils

import "encoding/json"

// An InterfaceMap is a wrapper object around map[string]interface{}.
// It's purpose is to more easily manipulate interface{} values in particular
// with multi-level data structures.
type InterfaceMap struct {
	d map[string]interface{}
}

// NewInterfaceMap creates and returns an InterfaceMap.
func NewInterfaceMap() *InterfaceMap {
	return &InterfaceMap{
		d: make(map[string]interface{}),
	}
}

// NewMap creates an InterfaceMap using an existing map.
func NewMap(data map[string]interface{}) *InterfaceMap {
	return &InterfaceMap{
		d: data,
	}
}

// Set the key to val.
func (m *InterfaceMap) Set(key string, val interface{}) {
	m.d[key] = val
}

// Get the value of key.
func (m *InterfaceMap) Get(key string) interface{} {
	return m.d[key]
}

// GetOK gets the value of key and if it exists.
func (m *InterfaceMap) GetOK(key string) (interface{}, bool) {
	v, ok := m.d[key]
	return v, ok
}

// KeyExists returns if the key exists in the map.
func (m *InterfaceMap) KeyExists(key string) bool {
	_, ok := m.d[key]
	return ok
}

// Delete key from the map.
func (m *InterfaceMap) Delete(key string) {
	delete(m.d, key)
}

// Len returns the count of items in the map.
func (m *InterfaceMap) Len() int {
	return len(m.d)
}

// Copy returns a new InterfaceMap with copies of the
// keys and values of m.
func (m *InterfaceMap) Copy() *InterfaceMap {
	n := NewInterfaceMap()
	for k, v := range m.d {
		n.d[k] = v
	}
	return n
}

// MarshalJSON marshals the underlying map instead of the struct.
func (m *InterfaceMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.d)
}

// UnmarshalJSON unmarshals the underlying map instead of the struct.
func (m *InterfaceMap) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, m.d)
}
