package event

import (
	"fmt"
	"testing"
	"time"
)

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

func TestEventToString(t *testing.T) {
	now := time.Now().UTC()
	e := New("Hello, world")
	e.SetTimestamp(now)

	s := e.String()
	expected := fmt.Sprintf("%s: Hello, world", now.Format(time.RFC3339))
	if s != expected {
		t.Errorf("Event to string. Expected %s, got %s", expected, s)
	}
}

func TestSprintf(t *testing.T) {
	now := time.Now()
	e := New("Hello, world")
	e.SetTimestamp(now)
	e.SetField("@source", "test")

	s := e.Sprintf("%{@timestamp} %{@source}: %{message}")
	expected := fmt.Sprintf("%s test: Hello, world", now.String())
	if s != expected {
		t.Errorf("Sprintf default time format. Expected %s, got %s", expected, s)
	}

	s = e.Sprintf("%{+2006-01-02T15:04:05Z07:00} %{@source}: %{message}")
	expected = fmt.Sprintf("%s test: Hello, world", now.UTC().Format(time.RFC3339))
	if s != expected {
		t.Errorf("Sprintf custom time format. Expected %s, got %s", expected, s)
	}
}

func TestSprintfNested(t *testing.T) {
	now := time.Now().UTC()
	e := New("Hello, world")
	e.SetTimestamp(now)
	e.SetField("field1.subfield1", "test")

	s := e.Sprintf("%{@timestamp}: %{message}: %{field1.subfield1}")
	expected := fmt.Sprintf("%s: Hello, world: test", now.String())
	if s != expected {
		t.Errorf("Sprintf Spartan nested. Expected %s, got %s", expected, s)
	}

	s = e.Sprintf("%{@timestamp}: %{message}: %{[field1][subfield1]}")
	expected = fmt.Sprintf("%s: Hello, world: test", now.String())
	if s != expected {
		t.Errorf("Sprintf Logstash nested. Expected %s, got %s", expected, s)
	}
}

func BenchmarkEventSprintf(b *testing.B) {
	now := time.Now()
	e := New("Hello, world")
	e.SetTimestamp(now)
	e.SetField("@source", "test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Sprintf("%{@timestamp} %{@source}: %{message}")
	}
}
