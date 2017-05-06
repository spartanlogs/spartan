package event

import (
	"fmt"
	"time"

	"github.com/spartanlogs/spartan/utils"
)

// Fields names for modules to use if needed
const (
	MessageField   = "message"
	TimestampField = "@timestamp"
	TypeField      = "type"
	TagsField      = "tags"
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
	data utils.InterfaceMap
}

// New creates a new Event object setting its message to message.
// @timestamp is set to time.Now() and tags and data are empty.
func New(message string) *Event {
	e := &Event{
		data: utils.NewInterfaceMap(),
	}
	e.SetTimestamp(time.Now())
	e.ResetTags()
	e.SetMessage(message)
	e.data.Set(TypeField, "")
	return e
}

func (e *Event) String() string {
	return fmt.Sprintf("%s: %s", e.GetTimestamp().Format(time.RFC3339), e.GetMessage())
}

// Data returns a copy of the underlying map with the events field data.
// The returned map is safe for the caller to manipulate as it's a copy
// of the underlying map in the Event. This method
// is intented for output/codecs modules to encode the Event.
func (e *Event) Data() utils.InterfaceMap {
	return e.data.Copy()
}

// SetField the field key to val.
func (e *Event) SetField(key string, val interface{}) {
	if key == TypeField {
		e.SetType(val.(string))
		return
	}
	e.data.Set(key, val)
}

// GetField the value of field key
func (e *Event) GetField(key string) interface{} {
	return e.data.Get(key)
}

// DeleteField deletes the field key
func (e *Event) DeleteField(key string) {
	e.data.Delete(key)
}

// HasField returns if the field key exists
func (e *Event) HasField(key string) bool {
	return e.data.HasKey(key)
}

// AddTag will add "tag" to the tags field if the tag doesn't
// already exist.
func (e *Event) AddTag(tag string) {
	if !e.HasTag(tag) {
		tags := e.GetField(TagsField).([]string)
		e.SetTags(append(tags, tag))
	}
}

// DeleteTag will delete "tag" from the tags field.
func (e *Event) DeleteTag(tag string) {
	i := e.hasTagIndex(tag)
	if i < 0 {
		return
	}
	tags := e.GetField(TagsField).([]string)
	e.SetTags(append(tags[:i], tags[i+1:]...))
}

// ResetTags removes all tags on the Event.
func (e *Event) ResetTags() {
	e.SetTags(make([]string, 0))
}

// HasTag returns if "tag" exists in the tags field.
func (e *Event) HasTag(tag string) bool {
	return e.hasTagIndex(tag) > -1
}

func (e *Event) hasTagIndex(tag string) int {
	for i, t := range e.GetField(TagsField).([]string) {
		if t == tag {
			return i
		}
	}
	return -1
}

// SetTags will replace the entire tag slice with the given slice.
func (e *Event) SetTags(tags []string) {
	e.SetField(TagsField, tags)
}

// GetTags returns the full tags field.
func (e *Event) GetTags() []string {
	return e.GetField(TagsField).([]string)
}

// SetMessage sets the Event message to s. Setting the message to
// an empty string will remove it from output.
func (e *Event) SetMessage(s string) {
	e.SetField(MessageField, s)
}

// GetMessage returns the current message.
func (e *Event) GetMessage() string {
	m := e.GetField(MessageField)
	if m == nil {
		return ""
	}
	return m.(string)
}

// SetType sets the type of the Event. Once type is set, it can't be changed.
func (e *Event) SetType(val string) {
	if e.GetType() == "" {
		e.data.Set(TypeField, val)
	}
}

// GetType returns the current Event type.
func (e *Event) GetType() string {
	v := e.GetField(TypeField)
	if v == nil {
		return ""
	}
	return v.(string)
}

// SetTimestamp sets the Event's canonical timestamp.
func (e *Event) SetTimestamp(t time.Time) {
	e.SetField(TimestampField, t)
}

// GetTimestamp returns the Events canonical timestamp.
func (e *Event) GetTimestamp() time.Time {
	return e.GetField(TimestampField).(time.Time)
}
