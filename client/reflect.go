package client

import (
	"reflect"
	"strings"
)

func appendDelim(s, postfix, delim string) string {
	if s != "" {
		s += delim
	}
	s += postfix
	return s
}

func keyToQuery(key interface{}) string {
	type serializable interface {
		Query() string
	}

	result := ""

	t := reflect.TypeOf(key)
	v := reflect.ValueOf(key)
	if t.Kind() != reflect.Struct {
		return result
	}

	nFields := t.NumField()
	for i := 0; i < nFields; i++ {
		f := t.Field(i)
		json := f.Tag.Get("json")
		if json != "" {
			json = strings.Split(json, ",")[0]
			val := v.Field(i)
			if inter, ok := val.Interface().(serializable); ok {
				result = appendDelim(result, json+"="+inter.Query(), "&")
			}
		}
	}
	return result
}

func getEntityName(entity interface{}) (string, error) {
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", ErrInvalidEntity
	}
	if t.NumField() < 1 {
		return "", ErrInvalidEntity
	}
	f := t.Field(0)
	if result, ok := f.Tag.Lookup("typename"); ok {
		return result, nil
	}
	return "", ErrInvalidEntity
}
