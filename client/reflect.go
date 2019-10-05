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

func paramsToQuery(params interface{}) string {
	type serializable interface {
		Query() string
	}

	result := ""

	t := reflect.TypeOf(params)
	v := reflect.ValueOf(params)
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
				result = appendDelim(result, json+"="+inter.Query(), ",")
			}
		}
	}
	return result
}

func getEntityName(entity interface{}) (string, error) {
	t := reflect.TypeOf(entity)
	for t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr {
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

func setClientToSlice(slice interface{}, c *Client) {
	t := reflect.TypeOf(slice)
	s := reflect.ValueOf(slice)
	for t.Kind() == reflect.Ptr {
		slice = s.Elem().Interface()
		t = reflect.TypeOf(slice)
		s = reflect.ValueOf(slice)
	}
	if t.Kind() != reflect.Slice {
		return
	}

	for i := 0; i < s.Len(); i++ {
		if e, ok := s.Index(i).Addr().Interface().(IEntity); ok {
			e.SetClient__(c)
		}
	}
}
