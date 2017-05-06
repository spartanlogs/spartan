package config

import (
	"testing"

	"github.com/spartanlogs/spartan/utils"
)

func TestNoSettings(t *testing.T) {
	m := utils.NewInterfaceMap()

	if err := VerifySettings(m, nil); err != nil {
		t.Errorf("Unexpected error %s", err)
	}
}

func TestSimpleSettings(t *testing.T) {
	m := utils.NewInterfaceMap()
	m.Set("field", "message")
	m.Set("action", "delete")

	settings := []Setting{
		{
			Name:     "field",
			Type:     String,
			Required: true,
		},
		{
			Name:     "action",
			Type:     String,
			Required: false,
			Default:  "reverse",
		},
	}

	if err := VerifySettings(m, settings); err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	m.Set("action", 42)

	if err := VerifySettings(m, settings); err == nil {
		t.Error("Error not received. Expected invalid setting type")
	}
}

func TestDefaultSettings(t *testing.T) {
	m := utils.NewInterfaceMap()
	m.Set("field", "message")

	settings := []Setting{
		{
			Name:     "field",
			Type:     String,
			Required: true,
		},
		{
			Name:    "action",
			Type:    String,
			Default: "reverse",
		},
	}

	VerifySettings(m, settings)

	v := m.Get("action").(string)
	if v != "reverse" {
		t.Errorf("Incorrect default value. Expected reverse, got %s", v)
	}
}

func TestArraySetting(t *testing.T) {
	m := utils.NewInterfaceMap()
	m.Set("fields", []string{"message", "type"})

	settings := []Setting{
		{
			Name:     "fields",
			Type:     Array,
			ElemType: &Setting{Type: String},
			Required: true,
		},
	}

	if err := VerifySettings(m, settings); err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	m.Set("fields", []int{42, 43})
	if err := VerifySettings(m, settings); err == nil {
		t.Error("Expected invalid element type error")
	}
}

func TestMapSetting(t *testing.T) {
	m := utils.NewInterfaceMap()
	m.Set("fields", utils.InterfaceMap{
		"message": true,
		"type":    false,
	})
	m.Set("actions", utils.InterfaceMap{
		"random":  "yep",
		"number?": 42,
	})

	settings := []Setting{
		{ // Test full verification of a map
			Name: "fields",
			Type: Map,
			MapDef: []Setting{
				{
					Name: "message",
					Type: Bool,
				},
				{
					Name: "type",
					Type: Bool,
				},
			},
			Required: true,
		},
		{ // Simple test that it is a map
			Name: "actions",
			Type: Map,
		},
	}

	if err := VerifySettings(m, settings); err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	m.Set("fields", utils.InterfaceMap{
		"message": true,
		"type":    42,
	})
	if err := VerifySettings(m, settings); err == nil {
		t.Error("Expected invalid element type error")
	}
}
