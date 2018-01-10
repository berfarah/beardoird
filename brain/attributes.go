package brain

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

const tagName = "attr"

type tagOptions struct {
	skip  bool // "-"
	name  string
	index bool
}

func getTagOpts(tag reflect.StructTag) tagOptions {
	t := tag.Get(tagName)

	if t == "" || t == "-" {
		return tagOptions{skip: true}
	}

	var opts tagOptions
	parts := strings.Split(t, ",")
	opts.name = parts[0]
	for _, s := range parts[1:] {
		switch s {
		case "index":
			opts.omitzero = true
		}
	}
	return opts
}

func (c *Client) Write(key string, s interface{}) error {
	v := reflect.ValueOf(s)
	fields := make(map[string]interface{})
	indexedFields := make(map[string]interface{})

	for i := 0; v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		opts := getTagOpts(fieldType.Tag)
		if opts.skip {
			continue
		}

		attributeName := fieldType.Name
		if opts.name != "" {
			attributeName = opts.name
		}

		if opts.index {
			indexedFields[attributeName] = field
		}

		fields[attributeName] = c.encode(field.Interface())
	}

	c.client.HMSet(key, fields)
	for name, value := range indexedFields {
		c.client.Set(key+":by:"+name, value, 0)
	}
}

func (c *Client) idOf(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	default:
		return "", errors.New("ID must be string or int")
	}
}

func (c *Client) encode(value reflect.Value) []byte {
}

func (c *Client) index(key, name string, value reflect.Value) {
	switch value.Interface().(type) {
	case reflect.Bool:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	}

	c.client.Set(key+":by:"+name, str, 0)
}

func Read(s interface{}) error {
}
