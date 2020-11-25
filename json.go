// Package jsonpathgo provides simple path based getters for accessing JSON primitives
package jsonpathgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type jsonValue interface {
	isJson()
}

type jsonNumber struct {
	Value float64
}

type jsonString struct {
	Value string
}

type jsonBool struct {
	Value bool
}

func (_ jsonNumber) isJson() {}
func (_ jsonString) isJson() {}
func (_ jsonBool) isJson()   {}

func interfaceForFirstKey(m map[string]interface{}) (interface{}, error) {
	for k, v := range m {
		switch v.(type) {
		case map[string]interface{}:
			return v, nil
		default:
			return nil, errors.New(fmt.Sprintf("* in path but matched first key '%s' wasn't a struct but %T", k, v))
		}
	}

	return nil, errors.New("No first key for empty struct")
}

func interfaceForArrayIndex(key string, index uint64, m map[string]interface{}) (interface{}, error) {
	v := m[key]

	switch vv := v.(type) {
	case []interface{}:
		if index >= uint64(len(vv)) {
			return nil, errors.New(fmt.Sprintf("index %d out of bounds for %s (had %d elements)", index, key, len(vv)))
		}
		return vv[index], nil
	default:
		return nil, errors.New(fmt.Sprintf("%s did not point to an array but %T", key, v))
	}
}

func jsonValueFactory(key string, v interface{}) (jsonValue, error) {
	switch vv := v.(type) {
	case string:
		return jsonString{vv}, nil
	case float64:
		return jsonNumber{vv}, nil
	case bool:
		return jsonBool{vv}, nil
	case interface{}: // xxx todo
		return nil, nil
	case nil:
		return nil, nil
	default:
		return nil, errors.New(fmt.Sprintf("%s has unexpected type %T", key, v))
	}
}

var arrayRe *regexp.Regexp = regexp.MustCompile("^([^[]*)\\[([0-9-]+)\\]")

func getIndexedKey(key string) (*string, *uint64, error) {
	if arrayRe.MatchString(key) {
		matches := arrayRe.FindStringSubmatch(key)
		index, err := strconv.ParseUint(matches[2], 0, 64)
		if err != nil {
			return nil, nil, err
		}
		return &matches[1], &index, nil
	}
	return nil, nil, nil
}

func headOfPath(path string) (head string, remainingPath string, err error) {
	headAndRestOfPath := strings.SplitN(path, ".", 2)

	if len(headAndRestOfPath) == 0 {
		return "", "", errors.New("Unexpected end of path")
	}

	head = headAndRestOfPath[0]

	if len(headAndRestOfPath) > 1 {
		remainingPath = headAndRestOfPath[1]
	} else {
		remainingPath = ""
	}

	return
}

func getJsonPathInterface(path string, m map[string]interface{}) (jsonValue, error) {
	key, remainingPath, err := headOfPath(path)

	if err != nil {
		return nil, err
	}

	indexKey, index, err := getIndexedKey(key)

	if err != nil {
		return nil, err
	}

	var v interface{}

	if key == "*" {
		v, err = interfaceForFirstKey(m)
	} else if indexKey != nil {
		v, err = interfaceForArrayIndex(*indexKey, *index, m)
	} else {
		v = m[key]
	}

	if err != nil {
		return nil, err
	}

	if len(remainingPath) == 0 {
		return jsonValueFactory(key, v)
	}

	switch v.(type) {
	case map[string]interface{}:
		m = v.(map[string]interface{})
	default:
		return nil, errors.New(fmt.Sprintf("Expected struct for %s, got %T", key, v))
	}

	return getJsonPathInterface(remainingPath, m)
}

func getJsonPath(path string, jsonInput []byte) (jsonValue, error) {
	var i interface{}

	err := json.Unmarshal(jsonInput, &i)
	if err != nil {
		return nil, err
	}

	return getJsonPathInterface(path, i.(map[string]interface{}))
}

// GetJsonPathString returns a string at a path from JSON.
// Given JSON
//     {"foo":{"123":{"bar":["41","42"]}}}
// path
//     foo.*.bar[1]
// would return "42"
func GetJsonPathString(path string, jsonInput []byte) (*string, error) {
	v, err := getJsonPath(path, jsonInput)

	if err != nil || v == nil {
		return nil, err
	}

	switch vv := v.(type) {
	case jsonString:
		return &vv.Value, nil
	default:
		return nil, errors.New(fmt.Sprintf("Expected path '%s' to point to string, got %T", path, v))
	}
}

// GetJsonPathString returns a bool at a path from JSON.
// Given JSON
//     {"foo":{"123":{"bar":[true,false]}}}
// path
//     foo.*.bar[1]
// would return false
func GetJsonPathBool(path string, jsonInput []byte) (*bool, error) {
	v, err := getJsonPath(path, jsonInput)

	if err != nil || v == nil {
		return nil, err
	}

	switch vv := v.(type) {
	case jsonBool:
		return &vv.Value, nil
	default:
		return nil, errors.New(fmt.Sprintf("Expected path '%s' to point to boolean, got %T", path, v))
	}
}

// GetJsonPathString returns a number at a path from JSON.
// Given JSON
//     {"foo":{"123":{"bar":[41,42]}}}
// path
//     foo.*.bar[1]
// would return 42
func GetJsonPathNumber(path string, jsonInput []byte) (*float64, error) {
	v, err := getJsonPath(path, jsonInput)

	if err != nil || v == nil {
		return nil, err
	}

	switch vv := v.(type) {
	case jsonNumber:
		return &vv.Value, nil
	default:
		return nil, errors.New(fmt.Sprintf("Expected path '%s' to point to number, got %T", path, v))
	}
}
