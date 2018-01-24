package main

import (
	"net/http"
	"os"

	"github.com/berfarah/beardroid/ludlow"
	"github.com/berfarah/gobot"
	"github.com/berfarah/gobot-slack"
	"github.com/berfarah/gobot-store-redis"
)

func main() {
	adapter := slack.New(os.Getenv("SLACK_TOKEN"))
	r := gobot.New(adapter)
	r.Install(
		redis.New(os.Getenv("REDIS_URL")),
		ludlow.Plugin,
	)

	r.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	}).Methods("GET")

	r.Debug(true)
	r.Run(":" + os.Getenv("PORT"))
}
