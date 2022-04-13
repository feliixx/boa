package boa

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

var (
	config   map[string]any
	defaults = map[string]any{}
)

// ParseConfig reads the config from an io.Reader.
func ParseConfig(jsonConfig io.Reader) error {

	d := json.NewDecoder(jsonConfig)
	d.UseNumber()

	return d.Decode(&config)
}

// SetDefault set the default value for this key.
func SetDefault(key string, value any) {

	switch reflect.TypeOf(value).Kind() {

	case reflect.Int,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint32,
		reflect.Uint64:
		defaults[key] = json.Number(fmt.Sprintf("%d", value))

	case reflect.Float64:
		defaults[key] = json.Number(strconv.FormatFloat(value.(float64), 'g', -1, 64))
	default:
		defaults[key] = value
	}

}

// GetString returns the value associated with the key as a string.
func GetString(key string) string {
	return cast[string](getValue(key))
}

// GetBool returns the value associated with the key as a bool.
func GetBool(key string) bool {
	return cast[bool](getValue(key))
}

// GetInt returns the value associated with the key as an int.
func GetInt(key string) int {

	n := cast[json.Number](getValue(key))
	v, err := strconv.ParseInt(string(n), 10, 64)
	if err != nil {
		panic("Can't parse number '%s' as an int")
	}
	return int(v)
}

// GetInt32 returns the value associated with the key as an int32
func GetInt32(key string) int32 {
	n := cast[json.Number](getValue(key))
	v, err := strconv.ParseInt(string(n), 10, 32)
	if err != nil {
		panic("Can't parse number '%s' as an int32")
	}
	return int32(v)
}

// GetInt64 returns the value associated with the key as an int64.
func GetInt64(key string) int64 {
	n := cast[json.Number](getValue(key))
	v, err := strconv.ParseInt(string(n), 10, 64)
	if err != nil {
		panic("Can't parse number '%s' as an int64")
	}
	return v
}

// GetUint returns the value associated with the key as an uint.
func GetUint(key string) uint {
	n := cast[json.Number](getValue(key))
	v, err := strconv.ParseUint(string(n), 10, 64)
	if err != nil {
		panic("Can't parse number '%s' as an uint")
	}
	return uint(v)
}

// GetUint32 returns the value associated with the key as an uint32.
func GetUint32(key string) uint32 {
	n := cast[json.Number](getValue(key))
	v, err := strconv.ParseUint(string(n), 10, 32)
	if err != nil {
		panic("Can't parse number '%s' as an uint32")
	}
	return uint32(v)
}

// GetUint64 returns the value associated with the key as an uint64.
func GetUint64(key string) uint64 {
	n := cast[json.Number](getValue(key))
	v, err := strconv.ParseUint(string(n), 10, 64)
	if err != nil {
		panic("Can't parse number '%s' as an uint64")
	}
	return uint64(v)
}

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key string) float64 {
	n := cast[json.Number](getValue(key))
	v, err := strconv.ParseFloat(string(n), 64)
	if err != nil {
		panic("Can't parse number '%s' as a float64")
	}
	return v
}

// GetAny returns any value associated with the key.
func GetAny(key string) any {
	return getValue(key)
}

func getValue(key string) any {

	nProp := strings.Split(key, ".")
	nested := config

	for i, prop := range nProp {

		if i == len(nProp)-1 {

			if v, ok := nested[prop]; ok {
				return v
			}

			path := "'" + strings.Join(nProp[:i], ".") + "'"
			if path == "''" {
				path = "root object"
			}

			return getDefaultValueOrPanic(
				key,
				fmt.Sprintf("%s has no key '%s'", path, prop),
			)
		}

		if _, ok := nested[prop].(map[string]any); !ok {
			return getDefaultValueOrPanic(
				key,
				fmt.Sprintf("%s is not an object", strings.Join(nProp[:i+1], ".")),
			)
		}

		nested = nested[prop].(map[string]any)
	}

	return nil
}

func getDefaultValueOrPanic(key, panicMsg string) any {
	v, ok := defaults[key]
	if !ok {
		panic(panicMsg)
	}
	return v
}

func cast[T any](v any) T {
	s, ok := v.(T)
	if !ok {
		panic(fmt.Sprintf("'%v' is not a %s", v, reflect.TypeOf(s)))
	}
	return s
}
