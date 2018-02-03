package main

import (
	"net/http"
	"os"

	"github.com/berfarah/beardroid/plugins/inventory/uniqlo"
	"github.com/botopolis/bot"
	"github.com/botopolis/redis"
	"github.com/botopolis/slack"
)

func main() {
	adapter := slack.New(os.Getenv("SLACK_TOKEN"))
	r := bot.New(
		adapter,
		redis.New(os.Getenv("REDIS_URL")),
		&uniqlo.Plugin{},
	)

	r.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	}).Methods("GET")

	r.Debug(true)
	r.Run()
}
