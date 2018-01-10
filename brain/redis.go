package brain

import "github.com/go-redis/redis"

type Client struct {
	client redis.Client
}

func New() *Client {
	return redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
		},
	)
}

func (c *Client) Set() {
}
