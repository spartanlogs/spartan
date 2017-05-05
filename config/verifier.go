package config

import (
	"errors"
	"fmt"

	"github.com/lfkeitel/spartan/utils"
)

// SettingType defines the type of a setting
type SettingType int

const (
	// Any skips type verification
	Any SettingType = iota
	// String setting
	String
	// Int setting (golang's int type)
	Int
	// Float setting (golang's float64 type)
	Float
	// Bool setting (true/false)
	Bool
	// Map setting
	Map
	// Array setting (elements can be Int, String, or Array)
	Array
)

// A Setting defines a single setting for a module
type Setting struct {
	// Name is used for the root or map key names
	Name     string
	Required bool

	// If value is nil, it will be assigned Default. Default is NOT verified
	// to be the right type for the setting. It's up to the module to ensure
	// the default value is valid.
	Default interface{}
	Type    SettingType

	// ElemType defines the type of elements in an array. The only setting that's
	// used is Type. Arrays can only consist of Strings, Ints, or another array of
	// the former.
	ElemType *Setting

	// MapDef defines the key:value pairs of a Map type setting.
	MapDef []Setting
}

// VerifySettings is used to check an InterfaceMap against a schema of Setting values.
func VerifySettings(data utils.InterfaceMap, settings []Setting) error {
	if len(settings) == 0 {
		return nil
	}

	for _, setting := range settings {
		v, exists := data.GetOK(setting.Name)
		if !exists && setting.Required {
			return fmt.Errorf("Expected required setting %s", setting.Name)
		}

		// If there's no value, no need to continue verification
		if v == nil {
			if setting.Default != nil {
				data.Set(setting.Name, setting.Default)
			} else {
				setTypeDefault(data, setting)
			}
			continue
		}

		switch setting.Type {
		case Any:
			continue
		case String:
			if _, ok := v.(string); !ok {
				return fmt.Errorf("Setting %s expected string, got %T", setting.Name, v)
			}
		case Int:
			if _, ok := v.(int); !ok {
				return fmt.Errorf("Setting %s expected int, got %T", setting.Name, v)
			}
		case Float:
			if _, ok := v.(float64); !ok {
				return fmt.Errorf("Setting %s expected float, got %T", setting.Name, v)
			}
		case Bool:
			if _, ok := v.(bool); !ok {
				return fmt.Errorf("Setting %s expected true/false, got %T", setting.Name, v)
			}
		case Array:
			if err := verifyArray(v, setting.ElemType); err != nil {
				return err
			}
		case Map:
			m, ok := v.(utils.InterfaceMap)
			if !ok {
				return fmt.Errorf("Setting %s expected map of values", setting.Name)
			}
			if err := VerifySettings(m, setting.MapDef); err != nil {
				return err
			}
		}
	}

	return nil
}

func verifyArray(data interface{}, elemType *Setting) error {
	if a, ok := data.([]interface{}); ok && len(a) == 0 {
		return nil
	}

	switch elemType.Type {
	case String:
		if _, ok := data.([]string); !ok {
			return errors.New("Expected array of string")
		}
	case Int:
		if _, ok := data.([]int); !ok {
			return errors.New("Expected array of int")
		}
	case Array:
		return verifyArray(data, elemType.ElemType)
	default:
		return errors.New("invalid array element type")
	}

	return nil
}

func setTypeDefault(data utils.InterfaceMap, setting Setting) {
	switch setting.Type {
	case String:
		data.Set(setting.Name, "")
	case Int:
		data.Set(setting.Name, 0)
	case Float:
		data.Set(setting.Name, 0.0)
	case Bool:
		data.Set(setting.Name, false)
	case Array:
		setArrayTypeDefault(data, setting.Name, setting.ElemType)
	case Map:
		data.Set(setting.Name, utils.NewInterfaceMap())
	}
}

func setArrayTypeDefault(data utils.InterfaceMap, name string, elemType *Setting) {
	switch elemType.Type {
	case String:
		data.Set(name, []string{})
	case Int:
		data.Set(name, []int{})
	case Array:
		setArrayTypeDefault(data, elemType.Name, elemType.ElemType)
	}
}
