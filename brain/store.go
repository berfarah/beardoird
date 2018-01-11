package brain

import (
	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

// Resource is a storable interface in the brain
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
// func (u User) Namespace() string {
//   return "app:user"
// }
//
// // optionally:
// func (u User) Indices() map[string]string {
//   return map[string]string{
//     "name": u.Name,
//   }
// }
type Resource interface {
	ID() string
	Namespace() string
}

type indexer interface {
	Resource
	Indices() map[string]string
}

// Brain is the store
type Brain struct{ cache *cache.Codec }

// New creates a new connection to Redis
func New() *Brain {
	return &Brain{
		cache: &cache.Codec{
			Redis: redis.NewClient(&redis.Options{
				Addr: "localhost:6379",
			}),
			Marshal: func(v interface{}) ([]byte, error) {
				return msgpack.Marshal(v)
			},
			Unmarshal: func(b []byte, v interface{}) error {
				return msgpack.Unmarshal(b, v)
			},
		},
	}
}
