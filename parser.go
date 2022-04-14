package boa

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	val, ok := getValue(key)
	if !ok {
		return ""
	}
	return cast[string](val)
}

// GetBool returns the value associated with the key as a bool.
func GetBool(key string) bool {

	val, ok := getValue(key)
	if !ok {
		return false
	}
	return cast[bool](val)
}

// GetInt returns the value associated with the key as an int.
func GetInt(key string) int {

	val, ok := getValue(key)
	if !ok {
		return 0
	}

	n := cast[json.Number](val)
	i, err := strconv.ParseInt(string(n), 10, 64)
	if err != nil {
		panic("Can't parse number '%s' as an int")
	}
	return int(i)
}

// GetInt32 returns the value associated with the key as an int32
func GetInt32(key string) int32 {

	val, ok := getValue(key)
	if !ok {
		return 0
	}

	n := cast[json.Number](val)
	i, err := strconv.ParseInt(string(n), 10, 32)
	if err != nil {
		panic("Can't parse number '%s' as an int32")
	}
	return int32(i)
}

// GetInt64 returns the value associated with the key as an int64.
func GetInt64(key string) int64 {

	val, ok := getValue(key)
	if !ok {
		return 0
	}

	n := cast[json.Number](val)
	i, err := strconv.ParseInt(string(n), 10, 64)
	if err != nil {
		panic("Can't parse number '%s' as an int64")
	}
	return i
}

// GetUint returns the value associated with the key as an uint.
func GetUint(key string) uint {

	val, ok := getValue(key)
	if !ok {
		return 0
	}

	n := cast[json.Number](val)
	i, err := strconv.ParseUint(string(n), 10, 64)
	if err != nil {
		panic("Can't parse number '%s' as an uint")
	}
	return uint(i)
}

// GetUint32 returns the value associated with the key as an uint32.
func GetUint32(key string) uint32 {

	val, ok := getValue(key)
	if !ok {
		return 0
	}

	n := cast[json.Number](val)
	i, err := strconv.ParseUint(string(n), 10, 32)
	if err != nil {
		panic("Can't parse number '%s' as an uint32")
	}
	return uint32(i)
}

// GetUint64 returns the value associated with the key as an uint64.
func GetUint64(key string) uint64 {

	val, ok := getValue(key)
	if !ok {
		return 0
	}

	n := cast[json.Number](val)
	i, err := strconv.ParseUint(string(n), 10, 64)
	if err != nil {
		panic("Can't parse number '%s' as an uint64")
	}
	return uint64(i)
}

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key string) float64 {

	val, ok := getValue(key)
	if !ok {
		return 0.0
	}

	n := cast[json.Number](val)
	f, err := strconv.ParseFloat(string(n), 64)
	if err != nil {
		panic("Can't parse number '%s' as a float64")
	}
	return f
}

// GetAny returns any value associated with the key.
func GetAny(key string) any {
	v, _ := getValue(key)
	return v
}

func getValue(key string) (val any, ok bool) {

	nProp := strings.Split(key, ".")
	nested := config

	for i, prop := range nProp {

		if i == len(nProp)-1 {

			if v, ok := nested[prop]; ok {
				return v, ok
			}
			return getDefault(key)
		}

		if _, ok := nested[prop].(map[string]any); !ok {
			return getDefault(key)
		}

		nested = nested[prop].(map[string]any)
	}
	return nil, false
}

func getDefault(key string) (val any, ok bool) {
	v, ok := defaults[key]
	if !ok {
		log.Printf("no value found for key '%s', using nil value instead", key)
	}
	return v, ok
}

func cast[T any](v any) T {
	s, ok := v.(T)
	if !ok {
		panic(fmt.Sprintf("'%v' is not a %s", v, reflect.TypeOf(s)))
	}
	return s
}
