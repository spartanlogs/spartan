package utils

import (
	"reflect"
	"testing"
)

func TestSimpleSetGet(t *testing.T) {
	m := NewInterfaceMap()

	m.Set("key1", "value1")
	m.Set("key2", "value2")

	v, ok := m.GetOK("key1")
	if !ok {
		t.Error("Value doesn't exist")
	}
	vs, ok := v.(string)
	if !ok {
		t.Error("Value is not a string")
	}
	if ok && vs != "value1" {
		t.Errorf("Expected value1, got %s", v)
	}

	v, ok = m.GetOK("key2")
	if !ok {
		t.Error("Value doesn't exist")
	}
	vs, ok = v.(string)
	if !ok {
		t.Error("Value is not a string")
	}
	if ok && vs != "value2" {
		t.Errorf("Expected value2, got %s", v)
	}

	m.Set("key1", "value3")

	v, ok = m.GetOK("key1")
	if !ok {
		t.Error("Value doesn't exist")
	}
	vs, ok = v.(string)
	if !ok {
		t.Error("Value is not a string")
	}
	if ok && vs != "value3" {
		t.Errorf("Expected value3, got %s", v)
	}
}

func TestNestedSetGet(t *testing.T) {
	m := NewInterfaceMap()

	m.Set("key1.subKey1", "value1")
	m.Set("key1.subKey2", "value2")

	v, ok := m.GetOK("key1.subKey1")
	if !ok {
		t.Error("Value doesn't exist")
	}
	vs, ok := v.(string)
	if !ok {
		t.Error("Value is not a string")
	}
	if ok && vs != "value1" {
		t.Errorf("Expected value1, got %s", v)
	}

	v, ok = m.GetOK("key1.subKey2")
	if !ok {
		t.Error("Value doesn't exist")
	}
	vs, ok = v.(string)
	if !ok {
		t.Error("Value is not a string")
	}
	if ok && vs != "value2" {
		t.Errorf("Expected value2, got %s", v)
	}

	m.Set("key1.subKey1", "value3")

	v, ok = m.GetOK("key1.subKey1")
	if !ok {
		t.Error("Value doesn't exist")
	}
	vs, ok = v.(string)
	if !ok {
		t.Error("Value is not a string")
	}
	if ok && vs != "value3" {
		t.Errorf("Expected value3, got %s", v)
	}
}

func TestHasKey(t *testing.T) {
	m := NewInterfaceMap()

	m.Set("key1.some.thing", "some value")
	if !m.HasKey("key1.some.thing") {
		t.Error("Expected key to exist key1.some.thing")
	}
	if m.HasKey("key1.some.thing1") {
		t.Error("Expected key not to exist key1.some.thing1")
	}
}

func TestDeleteKey(t *testing.T) {
	m := NewInterfaceMap()

	m.Set("key1.some.thing", "some value")
	if !m.HasKey("key1.some.thing") {
		t.Error("Expected key to exist key1.some.thing")
	}

	m.Delete("key1.some.thing")

	if m.HasKey("key1.some.thing") {
		t.Error("Expected key not to exist key1.some.thing")
	}
}

func TestMapLen(t *testing.T) {
	m := NewInterfaceMap()

	if m.Len() != 0 {
		t.Errorf("Expected map length to be 0, got %d", m.Len())
	}

	m.Set("key1.some.thing", "some value")
	if m.Len() != 1 {
		t.Errorf("Expected map length to be 1, got %d", m.Len())
	}
}

func TestMapCopy(t *testing.T) {
	m1 := NewInterfaceMap()

	m1.Set("key1.some.thing", "some value")
	m1.Set("key2.some", "some value")
	m1.Set("key1.some.thing1", "some value")
	m1.Set("key1.some.thing.another", "some value")
	m1.Set("key23.some.thing", "some value")
	m1.Set("key1.some1.thing", "some value")

	m2 := m1.Copy()

	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("Expected copied map to be equal. M1: %#v, M2: %#v", m1, m2)
	}
}

func TestOverwritingMapKey(t *testing.T) {
	m := NewInterfaceMap()

	m.Set("key1.some.thing", "some value")
	if _, ok := m.Get("key1.some").(InterfaceMap); !ok {
		t.Errorf("Expected key1.some to be Map, got %T", m.Get("key1.some"))
	}

	m.Set("key1.some", "some value")
	if _, ok := m.Get("key1.some").(string); !ok {
		t.Errorf("Expected key1.some to be string, got %T", m.Get("key1.some"))
	}
}
