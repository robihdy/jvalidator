package jvalidator

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/mitchellh/mapstructure"
)

var (
	// EmailPattern is a reguler expression for validating email
	EmailPattern = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// JValidator is a struct containing data (unmarshaled json data) and Errors (invalid json values).
type JValidator struct {
	data   map[string]interface{}
	Errors errors
}

// New is a function to generate a JValidator struct.
func New(jsonData []byte, v interface{}) (*JValidator, error) {
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(data, v)
	if err != nil {
		return nil, err
	}
	return &JValidator{
		data:   data,
		Errors: errors(map[string][]string{}),
	}, nil
}

func (j *JValidator) Required(names ...string) {
	for _, name := range names {
		_, ok := j.data[name]
		if !ok {
			j.Errors.Add(name, "Cannot be blank.")
			return
		}

		val, ok := j.data[name].(string)
		if ok && strings.TrimSpace(val) == "" {
			j.Errors.Add(name, "Cannot be blank.")
		}
	}
}

func (j *JValidator) String(names ...string) {
	for _, name := range names {
		_, ok := j.data[name].(string)
		if !ok {
			j.Errors.Add(name, "Must contain a string value.")
		}
	}
}

func (j *JValidator) Number(names ...string) {
	for _, name := range names {
		_, ok := j.data[name].(float64)
		if !ok {
			j.Errors.Add(name, "Must contain a number.")
		}
	}
}

func (j *JValidator) MaxChar(name string, d int) {
	val := j.data[name].(string)
	if val == "" {
		return
	}
	if utf8.RuneCountInString(val) > d {
		j.Errors.Add(name, fmt.Sprintf("Too long (maximum is %d characters).", d))
	}
}

func (j *JValidator) MinChar(name string, d int) {
	val := j.data[name].(string)
	if val == "" {
		return
	}
	if utf8.RuneCountInString(val) < d {
		j.Errors.Add(name, fmt.Sprintf("Too short (minimum is %d characters).", d))
	}
}

func (j *JValidator) MatchPattern(name string, pattern *regexp.Regexp) {
	val := j.data[name].(string)
	if val == "" {
		return
	}
	if !pattern.MatchString(val) {
		j.Errors.Add(name, "Invalid value.")
	}
}

func (j *JValidator) Email(name string) {
	val := j.data[name].(string)
	if val == "" {
		return
	}
	if !EmailPattern.MatchString(val) {
		j.Errors.Add(name, "Must contain a valid email address.")
	}
}

func (j *JValidator) IsValid() bool {
	return len(j.Errors) == 0
}

func (j *JValidator) ToJSON() ([]byte, error) {
	return json.Marshal(j.Errors)
}
