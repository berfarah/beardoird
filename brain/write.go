package brain

import (
	"errors"

	"github.com/go-redis/cache"
)

func (c *Brain) writeIndices(s Resource) error {
	i, ok := s.(indexer)
	if !ok {
		return nil
	}

	for name, value := range i.Indices() {
		if value == "" {
			continue
		}

		key := i.Namespace() + "_by_" + name + ":" + value
		item := cache.Item{Key: key, Object: s.ID()}
		if err := c.cache.Set(&item); err != nil {
			return err
		}
	}

	return nil
}

func (c *Brain) Write(s Resource) error {
	if s.ID() == "" || s.Namespace() == "" {
		return errors.New("ID() and Namespace() cannot be empty")
	}
	item := cache.Item{Key: s.Namespace() + ":" + s.ID(), Object: s}
	if err := c.cache.Set(&item); err != nil {
		return err
	}

	return c.writeIndices(s)
}
