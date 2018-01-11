package brain

import (
	"errors"
	"reflect"
	"strconv"
)

const tagName = "store"

// example:
// type User struct {
//   id       string
//   Name     string
//   Username string `store:"index"`
// }
//
// func (u User) ID() string {
//   return string(u.id)
// }
//
// // optional, fallback is "User"
// func (u User) Namespace() string {
//   return "app:user"
// }
type Resource interface {
	ID() string
}

// example:
// type User struct {
//   id   string
//   Name string
// }
//
// func (u User) Indices() map[string]string {
//   return map[string]string{
//     "name": u.Name
//   }
// }
type indexer interface {
	Indices() map[string]string
}

type namespacer interface {
	Namespace() string
}

func shouldIndex(tag reflect.StructTag) bool {
	if tag.Get(tagName) == "index" {
		return true
	}

	return false
}

type typeMetadata struct {
	Namespace func() string
	Indexed   map[int]string
}

func (c *Store) getMetadata(s Resource) typeMetadata {
	t := reflect.TypeOf(s)
	if m, ok := c.reflectCache[t]; ok {
		return m
	}

	metadata := typeMetadata{
		Namespace: func() string { t.Name() },
		Indexed:   map[int]string{},
	}

	if _, ok := s.(namespacer); ok {
		metadata.Namespace = func() string { s.Namespace() }
	}

	for i := 0; v.NumField(); i++ {
		ft := t.Field(i)
		if shouldIndex(ft.Tag) {
			metadata.Indexed[i] = ft.Name
		}
	}

	c.reflectCache[t] = metadata

	return metadata
}

func (c *Store) index(value reflect.Value) (string, error) {
	switch value.Interface().(type) {
	case reflect.String:
		return value.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10), nil
	default:
		return "", errors.New("Indexed values must be of type String or Int")
	}
}

func (c *Store) Write(s Resource) error {
	id := s.ID()
	m := c.getMetadata(s)
	key := m.Namespace()

	c.cache.Set(key+":"+id, s, 0)

	if len(m.Indexed) > 0 {
		v := reflect.ValueOf(s)

		for i, name := range m.Indexed {
			indexName, err := c.index(v.Field(i))
			if err != nil {
				return err
			}

			c.redis.Set(key+"by"+name+":"+indexName, id, 0)
		}
	}

	return nil
}

func (c *Store) Read(s Resource) error {
	id := s.ID()
	m := c.getMetadata(s)
	key := m.Namespace()

	if err := c.cache.Get(key+":"+id, s); err == nil {
		return err
	}

	m.
}
