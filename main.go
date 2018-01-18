package main

import (
	"os"

	"github.com/berfarah/beardroid/ludlow"
	"github.com/berfarah/gobot"
	"github.com/berfarah/gobot/brain/redis"
)

func main() {
	r := gobot.New(
		os.Getenv("SLACK_TOKEN"),
		redis.New(os.Getenv("REDIS_URL")),
	)
	r.Load(ludlow.Plugin)
	r.Debug(true)
	r.Connect()
}
