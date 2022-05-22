package boa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var (
	config        = map[string]any{}
	flattenConfig = map[string]any{}
	defaults      = map[string]any{}
)

// ParseConfig reads the config from an io.Reader.
//
// The config may be in JSON or in JSONC ( json with comment)
// Allowed format for comments are
//  * single line ( //... )
//  * multiligne ( /*...*/ )
func ParseConfig(jsoncConfig io.Reader) error {

	jsonc, err := io.ReadAll(jsoncConfig)
	if err != nil {
		return fmt.Errorf("fail to read from reader: %v", err)
	}

	cleanJson := removeComment(jsonc)

	d := json.NewDecoder(bytes.NewReader(cleanJson))
	d.UseNumber()

	err = d.Decode(&config)
	if err != nil {
		return fmt.Errorf("fail to parse JSON: %v", err)
	}

	flatten("", config, flattenConfig)

	return nil
}

func removeComment(src []byte) []byte {

	output := make([]byte, 0, len(src))

	var prev byte
	inString := false
	inSingleLineComment := false
	inMultiligneComment := false

	for _, b := range src {

		switch b {

		case '"':
			inString = !inString
		case '/':
			if !inString && prev == '/' {
				output = output[0 : len(output)-1]
				inSingleLineComment = true
			}

			if inMultiligneComment && prev == '*' {
				inMultiligneComment = false
				continue
			}

		case '*':
			if !inString && prev == '/' {
				output = output[0 : len(output)-1]
				inMultiligneComment = true
			}
		case '\n':
			inSingleLineComment = false
		}

		prev = b

		if !inSingleLineComment && !inMultiligneComment {
			output = append(output, b)
		}
	}
	return output
}

func flatten(prefix string, src map[string]any, dst map[string]any) {

	if len(prefix) > 0 {
		prefix += "."
	}

	for k, v := range src {

		switch child := v.(type) {
		case map[string]any:
			flatten(prefix+k, child, dst)
		default:
			dst[prefix+k] = v
		}
	}
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
	return getValue[string](key)
}

// GetBool returns the value associated with the key as a bool.
func GetBool(key string) bool {
	return getValue[bool](key)
}

// GetInt returns the value associated with the key as an int.
func GetInt(key string) int {

	val := getValue[json.Number](key)

	i, err := strconv.ParseInt(string(val), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Can't parse number '%s' as an int", val))
	}
	return int(i)
}

// GetInt32 returns the value associated with the key as an int32
func GetInt32(key string) int32 {

	val := getValue[json.Number](key)

	i, err := strconv.ParseInt(string(val), 10, 32)
	if err != nil {
		panic(fmt.Sprintf("Can't parse number '%s' as an int32", val))
	}
	return int32(i)
}

// GetInt64 returns the value associated with the key as an int64.
func GetInt64(key string) int64 {

	val := getValue[json.Number](key)

	i, err := strconv.ParseInt(string(val), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Can't parse number '%s' as an int64", val))
	}
	return i
}

// GetUint returns the value associated with the key as an uint.
func GetUint(key string) uint {

	val := getValue[json.Number](key)

	n := cast[json.Number](val)
	i, err := strconv.ParseUint(string(n), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Can't parse number '%s' as an uint", val))
	}
	return uint(i)
}

// GetUint32 returns the value associated with the key as an uint32.
func GetUint32(key string) uint32 {

	val := getValue[json.Number](key)

	i, err := strconv.ParseUint(string(val), 10, 32)
	if err != nil {
		panic(fmt.Sprintf("Can't parse number '%s' as an uint32", val))
	}
	return uint32(i)
}

// GetUint64 returns the value associated with the key as an uint64.
func GetUint64(key string) uint64 {

	val := getValue[json.Number](key)

	i, err := strconv.ParseUint(string(val), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Can't parse number '%s' as an uint64", val))
	}
	return uint64(i)
}

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key string) float64 {

	val := getValue[json.Number](key)

	f, err := strconv.ParseFloat(string(val), 64)
	if err != nil {
		panic(fmt.Sprintf("Can't parse number '%s' as an float64", val))
	}
	return f
}

// GetAny returns any value associated with the key.
func GetAny(key string) any {
	return getValue[any](key)
}

// GetMap returns the map associated with the key.
// Numbers will be of type json.Number
// returns nil if the key does not exist
func GetMap(key string) map[string]any {

	nProp := strings.Split(key, ".")
	nested := config

	for i, prop := range nProp {

		if i == len(nProp)-1 {

			if v, ok := nested[prop]; ok {
				return cast[map[string]any](v)
			}
			continue
		}
		nested, _ = nested[prop].(map[string]any)
	}
	return getDefault[map[string]any](key)
}

func getValue[T any](key string) (val T) {

	v, ok := flattenConfig[key]
	if ok {
		return cast[T](v)
	}
	return getDefault[T](key)
}

func getDefault[T any](key string) (val T) {

	v, ok := defaults[key]
	if ok {
		return cast[T](v)
	}

	log.Printf("no value found for key '%s', using nil value instead", key)

	var zeroVal T
	if reflect.TypeOf(zeroVal) == reflect.TypeOf(json.Number("")) {
		return cast[T](json.Number("0"))
	}
	return zeroVal
}

func cast[T any](v any) T {

	s, ok := v.(T)
	if !ok {
		panic(fmt.Sprintf("'%v' is not a %s", v, reflect.TypeOf(s)))
	}
	return s
}
