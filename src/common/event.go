package common

import (
	"sync"
	"time"
)

type Event struct {
	sync.RWMutex
	timestamp time.Time
	etype     string
	message   string
	tags      []string
	data      map[string]interface{}
}

func NewEvent(message string) *Event {
	return &Event{
		timestamp: time.Now(),
		message:   message,
		tags:      make([]string, 0),
		data:      make(map[string]interface{}),
	}
}

func (e *Event) Squash() map[string]interface{} {
	e.Lock()
	if e.message != "" {
		e.data["message"] = e.message
	}
	e.data["type"] = e.etype
	e.data["@timestamp"] = e.timestamp
	e.data["tags"] = e.tags
	e.Unlock()
	return e.data
}

func (e *Event) Set(key string, val interface{}) {
	if e.setSpecial(key, val) {
		return
	}
	e.Lock()
	e.data[key] = val
	e.Unlock()
}

func (e *Event) setSpecial(key string, val interface{}) bool {
	switch key {
	case "message":
		if s, ok := val.(string); ok {
			e.SetMessage(s)
		}
		return true
	case "type":
		if s, ok := val.(string); ok {
			e.SetType(s)
		}
		return true
	case "@timestamp":
		if t, ok := val.(time.Time); ok {
			e.SetTimestamp(t)
		}
		return true
	}
	return false
}

func (e *Event) Get(key string) interface{} {
	if val, exists := e.getSpecial(key); exists {
		return val
	}
	e.RLock()
	defer e.RUnlock()
	if val, exists := e.data[key]; exists {
		return val
	}
	return nil
}

func (e *Event) getSpecial(key string) (interface{}, bool) {
	switch key {
	case "message":
		return e.message, true
	case "type":
		return e.etype, true
	case "tags":
		return e.tags, true
	case "@timestamp":
		return e.timestamp, true
	}
	return nil, false
}

func (e *Event) RemoveField(key string) {
	if ok := e.removeSpecial(key); ok {
		return
	}
	e.RLock()
	defer e.RUnlock()
	delete(e.data, key)
	return
}

func (e *Event) removeSpecial(key string) bool {
	switch key {
	case "message":
		e.message = ""
		return true
	case "tags":
		e.ResetTags()
		return true
	// @timestamp and type are protected fields,
	// they can not be removed.
	case "@timestamp":
		fallthrough
	case "type":
		return true
	}
	return false
}

func (e *Event) AddTag(tag string) {
	if !e.HasTag(tag) {
		e.tags = append(e.tags, tag)
	}
}

func (e *Event) RemoveTag(tag string) {
	i := e.hasTagIndex(tag)
	if i < 0 {
		return
	}
	e.tags = append(e.tags[:i], e.tags[i+1:]...)
}

func (e *Event) ResetTags() {
	e.tags = make([]string, 0)
}

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

func (e *Event) SetTags(tags []string) {
	e.tags = tags
}

func (e *Event) GetTags() []string {
	return e.tags
}

func (e *Event) SetMessage(s string) {
	e.message = s
}

func (e *Event) GetMessage() string {
	return e.message
}

func (e *Event) SetType(val string) {
	if e.etype == "" {
		e.etype = val
	}
}

func (e *Event) GetType() string {
	return e.etype
}

func (e *Event) SetTimestamp(t time.Time) {
	e.timestamp = t
}

func (e *Event) GetTimestamp() time.Time {
	return e.timestamp
}
