package main

import (
	"net/http"
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
	r.Install(
		ludlow.Plugin,
	)

	r.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	}).Methods("GET")

	r.Debug(true)
	r.Start(":" + os.Getenv("PORT"))
}
