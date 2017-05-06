package parser

import (
	"reflect"
	"testing"

	"github.com/spartanlogs/spartan/config/lexer"
	"github.com/spartanlogs/spartan/utils"
)

func TestEmptyArrayParser(t *testing.T) {
	l := lexer.NewString(`[]`)
	p := newParser(l)
	a, err := p.parseArray()
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	asSlice, ok := a.([]interface{})
	if !ok {
		t.Fatalf("Parsed array is not a []interface{}. Type: %s", reflect.TypeOf(a))
	}

	if len(asSlice) != 0 {
		t.Fatalf("Unexpected length. Expected 0, got %d", len(asSlice))
	}
}

func TestStringArrayParser(t *testing.T) {
	l := lexer.NewString(`["item1", "item2"]`)
	p := newParser(l)
	a, err := p.parseArray()
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	asSlice, ok := a.([]string)
	if !ok {
		t.Fatalf("Parsed array is not a []string. Type: %s", reflect.TypeOf(a))
	}

	if len(asSlice) != 2 {
		t.Fatalf("Unexpected length. Expected 2, got %d", len(asSlice))
	}

	if asSlice[0] != "item1" {
		t.Fatalf("Item 0 is %s, expected item1", asSlice[0])
	}

	if asSlice[1] != "item2" {
		t.Fatalf("Item 1 is %s, expected item2", asSlice[0])
	}
}

func TestIntArrayParser(t *testing.T) {
	l := lexer.NewString(`[123, 456]`)
	p := newParser(l)
	a, err := p.parseArray()
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	asSlice, ok := a.([]int)
	if !ok {
		t.Fatalf("Parsed array is not a []int. Type: %s", reflect.TypeOf(a))
	}

	if len(asSlice) != 2 {
		t.Fatalf("Unexpected length. Expected 2, got %d", len(asSlice))
	}

	if asSlice[0] != 123 {
		t.Fatalf("Item 0 is %d, expected 123", asSlice[0])
	}

	if asSlice[1] != 456 {
		t.Fatalf("Item 1 is %d, expected 456", asSlice[0])
	}
}

func TestSimpleMapParser(t *testing.T) {
	l := lexer.NewString(`{"key1" => "val1", "key2" => 789, "key3" => 5.65}`)
	p := newParser(l)
	m, err := p.parseMap()
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	expected := map[string]interface{}{
		"key1": "val1",
		"key2": 789,
		"key3": 5.65,
	}

	for k, v := range expected {
		val, exists := m.GetOK(k)
		if !exists {
			t.Fatalf("Key %s doesn't exist but was expected.", k)
		}

		switch vTyped := v.(type) {
		case string:
			valTyped, ok := val.(string)
			if !ok {
				t.Fatalf("Value expected to be string, got %s", reflect.TypeOf(val))
			}
			if valTyped != vTyped {
				t.Fatalf("Mismatched values for key %s. Expected %s, got %s",
					k, vTyped, valTyped)
			}
		case int:
			valTyped, ok := val.(int)
			if !ok {
				t.Fatalf("Value expected to be int, got %s", reflect.TypeOf(val))
			}
			if valTyped != vTyped {
				t.Fatalf("Mismatched values for key %s. Expected %d, got %d",
					k, vTyped, valTyped)
			}
		case float64:
			valTyped, ok := val.(float64)
			if !ok {
				t.Fatalf("Value expected to be float64, got %s", reflect.TypeOf(val))
			}
			if valTyped != vTyped {
				t.Fatalf("Mismatched values for key %s. Expected %g, got %g",
					k, vTyped, valTyped)
			}
		}
	}
}

func TestComplexMapParser(t *testing.T) {
	l := lexer.NewString(`{"key1" => ["val1", "val2"], "key2" => {"subkey1" => "otherval1"}}`)
	p := newParser(l)
	m, err := p.parseMap()
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	expected := utils.NewMap(map[string]interface{}{
		"key1": []string{"val1", "val2"},
		"key2": utils.NewMap(map[string]interface{}{"subkey1": "otherval1"}),
	})

	if !reflect.DeepEqual(expected, m) {
		t.Fatalf("maps not equal. Expected %#v,\n\ngot %#v", expected, m)
	}
}

func TestPipelineParser(t *testing.T) {
	l := lexer.NewString(`{
	date {
		patterns => ["YYYY-mm-dd HH:ss", "UNIX"]
		timezone => "America/Chicago"
	}

	mutate {
		add_fields => {
			"field1" => "value1"
		}
	}
}`)
	p := newParser(l)
	pipeline, err := p.parsePipelineDefs()
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	if len(pipeline) != 2 {
		t.Fatalf("Incorrect pipeline len. Expected 2, got %d", len(pipeline))
	}

	// Check date module
	date := pipeline[0]
	if date.Module != "date" {
		t.Fatalf("Incorrect module name. Expected date, got %s", date.Module)
	}

	dateOptions := utils.NewMap(map[string]interface{}{
		"patterns": []string{"YYYY-mm-dd HH:ss", "UNIX"},
		"timezone": "America/Chicago",
	})
	if !reflect.DeepEqual(dateOptions, date.Options) {
		t.Fatalf("Incorrect module options. Expected %#v\n\ngot %#v", dateOptions, date.Options)
	}
	if len(date.Connections) != 1 {
		t.Fatalf("Incorrect number of connections. Expected 1, got %d", date.Connections)
	}
	if date.Connections[0] != 1 {
		t.Fatalf("Incorrect connection. Expected index 1, got index %d", date.Connections[0])
	}

	// Check mutate module
	mutate := pipeline[1]
	if mutate.Module != "mutate" {
		t.Fatalf("Incorrect module name. Expected mutate, got %s", mutate.Module)
	}

	mutateOptions := utils.NewMap(map[string]interface{}{
		"add_fields": utils.NewMap(map[string]interface{}{
			"field1": "value1",
		}),
	})
	if !reflect.DeepEqual(mutateOptions, mutate.Options) {
		t.Fatalf("Incorrect module options. Expected %#v\n\ngot %#v", mutateOptions, mutate.Options)
	}
	if len(mutate.Connections) != 0 {
		t.Fatalf("Incorrect number of connections. Expected 0, got %d", mutate.Connections)
	}
}

var fullExpectedFile = &ParsedFile{
	Inputs: []*InputDef{{
		Module: "file",
		Options: utils.NewMap(map[string]interface{}{
			"path": "",
		}),
	}},

	Outputs: []*PipelineDef{{
		Module: "stdout",
		Options: utils.NewMap(map[string]interface{}{
			"codec": "json",
		}),
	}},

	Filters: []*PipelineDef{{
		Module: "grok",
		Options: utils.NewMap(map[string]interface{}{
			"field": "message",
			"patterns": []string{
				`^(?<logdate>%{MONTHDAY}[-]%{MONTH}[-]%{YEAR} %{TIME}) client %{IP:clientip}#%{POSINT:clientport} \(%{GREEDYDATA:query}\): query: %{GREEDYDATA:target} IN %{GREEDYDATA:querytype} \(%{IP:dns}\)$`,
			},
		}),
		Connections: []int{1},
	}, {
		Module: "date",
		Options: utils.NewMap(map[string]interface{}{
			"field":    "logdate",
			"patterns": []string{"dd-MMM-yyyy HH:mm:ss.SSS"},
			"timezone": "America/Chicago",
		}),
		Connections: []int{2},
	}, {
		Module: "mutate",
		Options: utils.NewMap(map[string]interface{}{
			"action": "remove_field",
			"fields": []string{"logdate", "message"},
		}),
	}},
}

func TestFileParser(t *testing.T) {
	pf, err := ParseFile("./testdata/filters.conf")
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	if !reflect.DeepEqual(fullExpectedFile, pf) {
		t.Fatal("Incorrect ParsedFile")
	}
}

func TestFileGlobParser(t *testing.T) {
	pf, err := ParseGlob("./testdata/filters/*.conf")
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	if !reflect.DeepEqual(fullExpectedFile, pf) {
		t.Fatal("Incorrect ParsedFile")
	}
}
