package common

import (
	"sync"
	"time"
)

const (
	messageField   = "message"
	timestampField = "@timestamp"
	typeField      = "type"
	tagsField      = "tags"
)

// An Event is the primary data structure passed around the system.
// An input creates an Event which is sent to the filter pipeline and
// eventually out to the output pipeline.
//
// The fields @timestamp, message, type, and tags are protected
// and should be manipulated using the corresponding Get and Set methods.
// If any of these fields are manipulated using the generic Get/Set methods,
// the request will be directed to the correct Get/Set method.
type Event struct {
	sync.RWMutex
	timestamp time.Time
	etype     string
	message   string
	tags      []string
	data      *InterfaceMap
}

// NewEvent creates a new Event object setting its message to message.
// @timestamp is set to time.Now() and tags and data are empty.
func NewEvent(message string) *Event {
	return &Event{
		timestamp: time.Now(),
		message:   message,
		tags:      make([]string, 0),
		data:      NewInterfaceMap(),
	}
}

// Squash reduces the Event to an InterfaceMap where the map keys
// are the Event's field names. The returned map is safe for the caller to
// manipulate as it's a copy of the underlying map in the Event. This method
// is intented for output/codecs modules to encode the Event.
func (e *Event) Squash() *InterfaceMap {
	e.Lock()
	dataCopy := e.data.Copy()

	if e.message != "" {
		dataCopy.Set(messageField, e.message)
	}
	dataCopy.Set(typeField, e.etype)
	dataCopy.Set(timestampField, e.timestamp)
	dataCopy.Set(tagsField, e.tags)
	e.Unlock()
	return dataCopy
}

// Set the field key to val.
func (e *Event) Set(key string, val interface{}) {
	if e.setSpecial(key, val) {
		return
	}
	e.Lock()
	e.data.Set(key, val)
	e.Unlock()
}

func (e *Event) setSpecial(key string, val interface{}) bool {
	switch key {
	case messageField:
		if s, ok := val.(string); ok {
			e.SetMessage(s)
		}
		return true
	case typeField:
		if s, ok := val.(string); ok {
			e.SetType(s)
		}
		return true
	case timestampField:
		if t, ok := val.(time.Time); ok {
			e.SetTimestamp(t)
		}
		return true
	}
	return false
}

// Get the value of field key
func (e *Event) Get(key string) interface{} {
	if val, exists := e.getSpecial(key); exists {
		return val
	}
	e.RLock()
	defer e.RUnlock()
	if val, exists := e.data.GetOK(key); exists {
		return val
	}
	return nil
}

func (e *Event) getSpecial(key string) (interface{}, bool) {
	switch key {
	case messageField:
		return e.GetMessage(), true
	case typeField:
		return e.GetType(), true
	case tagsField:
		return e.GetTags, true
	case timestampField:
		return e.GetTimestamp(), true
	}
	return nil, false
}

// RemoveField deletes the field key
func (e *Event) RemoveField(key string) {
	if ok := e.removeSpecial(key); ok {
		return
	}
	e.RLock()
	defer e.RUnlock()
	e.data.Delete(key)
	return
}

func (e *Event) removeSpecial(key string) bool {
	switch key {
	case messageField:
		e.message = ""
		return true
	case tagsField:
		e.ResetTags()
		return true
	// @timestamp and type are protected fields,
	// they can not be removed.
	case timestampField:
		fallthrough
	case typeField:
		return true
	}
	return false
}

// AddTag will add "tag" to the tags field if the tag doesn't
// already exist.
func (e *Event) AddTag(tag string) {
	if !e.HasTag(tag) {
		e.tags = append(e.tags, tag)
	}
}

// RemoveTag will delete "tag" from the tags field.
func (e *Event) RemoveTag(tag string) {
	i := e.hasTagIndex(tag)
	if i < 0 {
		return
	}
	e.tags = append(e.tags[:i], e.tags[i+1:]...)
}

// ResetTags removes all tags on the Event.
func (e *Event) ResetTags() {
	e.tags = make([]string, 0)
}

// HasTag returns if "tag" exists in the tags field.
func (e *Event) HasTag(tag string) bool {
	return e.hasTagIndex(tag) > -1
}

func (e *Event) hasTagIndex(tag string) int {
	for i, t := range e.tags {
		if t == tag {
			return i
		}
	}
	return -1
}

// SetTags will replace the entire tag slice with the given slice.
func (e *Event) SetTags(tags []string) {
	e.tags = tags
}

// GetTags returns the full tags field.
func (e *Event) GetTags() []string {
	return e.tags
}

// SetMessage sets the Event message to s. Setting the message to
// an empty string will remove it from output.
func (e *Event) SetMessage(s string) {
	e.message = s
}

// GetMessage returns the current message.
func (e *Event) GetMessage() string {
	return e.message
}

// SetType sets the type of the Event. Once type is set, it can't be changed.
func (e *Event) SetType(val string) {
	if e.etype == "" {
		e.etype = val
	}
}

// GetType returns the current Event type.
func (e *Event) GetType() string {
	return e.etype
}

// SetTimestamp sets the Event's canonical timestamp.
func (e *Event) SetTimestamp(t time.Time) {
	e.timestamp = t
}

// GetTimestamp returns the Events canonical timestamp.
func (e *Event) GetTimestamp() time.Time {
	return e.timestamp
}
