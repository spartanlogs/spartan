package event

import "testing"

func TestRemoveField(t *testing.T) {
	e := New("")
	e.SetField("field1", "value1")

	if !e.HasField("field1") {
		t.Error("Expected field1 to exist")
	}

	e.DeleteField("field1")

	if e.HasField("field1") {
		t.Error("Expected field1 not to exist")
	}

	if e.Data().HasKey("field1") {
		t.Error("Expected field1 not to exist")
	}
}
