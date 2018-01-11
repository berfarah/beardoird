package brain

import (
	"reflect"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

type reflectCache map[reflect.Type]typeMetadata

type Store struct {
	redis        redis.Client
	cache        *cache.Codec
	reflectCache map[reflect.Type]typeMetadata
}

func New() *Store {
	instance := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	codec := &cache.Codec{
		Redis: instance,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}

	return &Store{
		redis:        instance,
		cache:        codec,
		reflectCache: reflectCache{},
	}
}
