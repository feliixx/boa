package boa_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/feliixx/boa"
)

func TestGetString(t *testing.T) {

	getStringTests := []struct {
		name     string
		config   string
		property string
		expected string
		panic    bool
		panicMsg string
	}{
		{
			name:     "single key",
			config:   `{ "key": "value" }`,
			property: "key",
			expected: "value",
		},
		{
			name:     "non existing key",
			config:   `{ "key": "value" }`,
			property: "yellow",
			expected: "",
		},
		{
			name:     "one nested level",
			config:   `{ "root": { "key": "value"}}`,
			property: "root.key",
			expected: "value",
		},
		{
			name:     "non existing nested",
			config:   `{"root": {"first": {"key":"value"}}}`,
			property: "root.wrong.key",
			expected: "",
		},
		{
			name:     "non existing nested key",
			config:   `{"root": {"first": {"key":"value"}}}`,
			property: "root.first.yellow",
			expected: "",
		},
		{
			name:     "wrong type",
			config:   `{"root": {"first": {"key":true}}}`,
			property: "root.first.key",
			panic:    true,
			panicMsg: `'true' is not a string`,
		},
	}

	for _, test := range getStringTests {

		tt := test
		t.Run(tt.name, func(t *testing.T) {

			if tt.panic {
				defer func() {
					panic := recover()
					if panic == nil {
						t.Error("test should have triggered a panic")
					}

					if panic != tt.panicMsg {
						t.Errorf("expected '%s' but got '%s'", tt.panicMsg, panic)
					}
				}()
			}

			loadConfig(t, tt.config)

			got := boa.GetString(tt.property)
			if got != tt.expected {
				t.Errorf("expected '%s' but got '%s'", tt.expected, got)
			}
		})
	}
}

func TestAll(t *testing.T) {

	config := `
	{
		"root": {
			"string": "s",
			"boolTrue": true,
			"boolFalse": false,
			"int": -1,
			"int32": -32,
			"int64": -64,
			"uint": 1,
			"uint32": 32,
			"uint64": 64,
			"float64": 1.23842323
		}
	}`

	loadConfig(t, config)

	tests := []struct {
		name string
		want any
		got  any
	}{
		{
			name: "string",
			want: "s",
			got:  boa.GetString("root.string"),
		},
		{
			name: "bool true",
			want: true,
			got:  boa.GetBool("root.boolTrue"),
		}, {
			name: "bool false",
			want: false,
			got:  boa.GetBool("root.boolFalse"),
		},
		{
			name: "int",
			want: int(-1),
			got:  boa.GetInt("root.int"),
		},
		{
			name: "int32",
			want: int32(-32),
			got:  boa.GetInt32("root.int32"),
		},
		{
			name: "int64",
			want: int64(-64),
			got:  boa.GetInt64("root.int64"),
		},
		{
			name: "uint",
			want: uint(1),
			got:  boa.GetUint("root.uint"),
		},
		{
			name: "uint32",
			want: uint32(32),
			got:  boa.GetUint32("root.uint32"),
		},
		{
			name: "uint64",
			want: uint64(64),
			got:  boa.GetUint64("root.uint64"),
		},
		{
			name: "float64",
			want: 1.23842323,
			got:  boa.GetFloat64("root.float64"),
		},

		{
			name: "bool nil value",
			want: false,
			got:  boa.GetBool("nilvalue.bool"),
		},
		{
			name: "int nil value",
			want: int(0),
			got:  boa.GetInt("nilvalue.int"),
		},
		{
			name: "int32 nil value",
			want: int32(0),
			got:  boa.GetInt32("nilvalue.int32"),
		},
		{
			name: "int64 nil value",
			want: int64(0),
			got:  boa.GetInt64("nilvalue.int64"),
		},
		{
			name: "uint nil value",
			want: uint(0),
			got:  boa.GetUint("nilvalue.uint"),
		},
		{
			name: "uint32 nil value",
			want: uint32(0),
			got:  boa.GetUint32("nilvalue.uint32"),
		},
		{
			name: "uint64 nil value",
			want: uint64(0),
			got:  boa.GetUint64("nilvalue.uint64"),
		},
		{
			name: "float64 nil value",
			want: 0.0,
			got:  boa.GetFloat64("nilvalue.float64"),
		},
	}

	for _, test := range tests {

		tt := test
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != tt.got {
				t.Errorf("expected '%v' but got '%v'", tt.want, tt.got)
			}
		})
	}
}

func TestSetDefaults(t *testing.T) {

	config := `
	{
		"root": {
			"string": "s"
		}
	}`

	loadConfig(t, config)

	boa.SetDefault("root.string", "unused")
	boa.SetDefault("root.first", "first")
	boa.SetDefault("root.first.second", "second")
	boa.SetDefault("int", -12)
	boa.SetDefault("int32", -2334)
	boa.SetDefault("int64", -17286145274665)
	boa.SetDefault("uint", 12)
	boa.SetDefault("uint32", 2334)
	boa.SetDefault("uint64", 17286145274665)
	boa.SetDefault("float64", 0.72631524721)
	boa.SetDefault("precision.float", 1563246315263.35152323132)

	tests := []struct {
		name string
		want any
		got  any
	}{
		{
			name: "default not used if key exist",
			want: "s",
			got:  boa.GetString("root.string"),
		},
		{
			name: "default used if key does not exist",
			want: "first",
			got:  boa.GetString("root.first"),
		},
		{
			name: "default used if nested object doesn't exist",
			want: "second",
			got:  boa.GetString("root.first.second"),
		},
		{
			name: "default int",
			want: -12,
			got:  boa.GetInt("int"),
		},
		{
			name: "default int32",
			want: int32(-2334),
			got:  boa.GetInt32("int32"),
		},
		{
			name: "default int64",
			want: int64(-17286145274665),
			got:  boa.GetInt64("int64"),
		},
		{
			name: "default uint",
			want: uint(12),
			got:  boa.GetUint("uint"),
		},
		{
			name: "default uint32",
			want: uint32(2334),
			got:  boa.GetUint32("uint32"),
		},
		{
			name: "default uint64",
			want: uint64(17286145274665),
			got:  boa.GetUint64("uint64"),
		},
		{
			name: "default float",
			want: 0.72631524721,
			got:  boa.GetFloat64("float64"),
		},
		{
			name: "precision float",
			want: 1563246315263.35152323132,
			got:  boa.GetFloat64("precision.float"),
		},
	}

	for _, test := range tests {

		tt := test
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != tt.got {
				t.Errorf("expected '%v' but got '%v'", tt.want, tt.got)
			}
		})
	}
}

func TestRemoveComment(t *testing.T) {

	config := `
	{
		/* some 
		multiline 
		comment */
		"smtp": {
			// single line 
			"enabled": true, // with trailing space
			"host": "http://127.0.0.1",// without trailing space
			"port": 55,
			"pwd": "fhd/*|,;,bdo*/"
		} 			/**/


	}
	// `

	loadConfig(t, config)

	tests := []struct {
		name string
		want any
		got  any
	}{
		{
			name: "smtp.enabled",
			want: true,
			got:  boa.GetBool("smtp.enabled"),
		},
		{
			name: "smtp.host",
			want: "http://127.0.0.1",
			got:  boa.GetString("smtp.host"),
		},
		{
			name: "smtp.port",
			want: 55,
			got:  boa.GetInt("smtp.port"),
		},
		{
			name: "smtp.pwd",
			want: "fhd/*|,;,bdo*/",
			got:  boa.GetString("smtp.pwd"),
		},
	}

	for _, test := range tests {

		tt := test
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != tt.got {
				t.Errorf("expected '%v' but got '%v'", tt.want, tt.got)
			}
		})
	}
}

func TestGetMap(t *testing.T) {

	config := `{
		"root": {
			"object": {
				"key1": "val",
				"key2": true,
				"key3": 1
			}
		}
	}`

	loadConfig(t, config)

	result := boa.GetMap("root.object")
	if want, got := "val", result["key1"]; want != got {
		t.Errorf("expected '%v' but got '%v'", want, got)
	}
	if want, got := true, result["key2"]; want != got {
		t.Errorf("expected '%v' but got '%v'", want, got)
	}
	if want, got := json.Number("1"), result["key3"]; want != got {
		t.Errorf("expected '%v' but got '%v'", want, got)
	}

	nonExistingMap := boa.GetMap("non_existing")
	if nonExistingMap != nil {
		t.Errorf("non existing key should return nil, but got '%v'", nonExistingMap)
	}

	defer func(t *testing.T) {
		got := recover()
		if got == nil {
			t.Error("val not of type map should trigger a panic")
		}
		want := "'val' is not a map[string]interface {}"
		if want != got {
			t.Errorf("expected '%s' but got '%s'", want, got)
		}
	}(t)
	_ = boa.GetMap("root.object.key1")

	boa.SetDefault("non_existing", map[string]any{"k": "v"})
	nonExistingMap = boa.GetMap("non_existing")
	if len(nonExistingMap) != 1 {
		t.Errorf("after default set, non_existing should be a map of length 1, but got '%v'", nonExistingMap)
	}
}

func loadConfig(t *testing.T, config string) {
	err := boa.ParseConfig(strings.NewReader(config))
	if err != nil {
		t.Error(err)
	}
}
