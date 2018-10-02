package main

import (
	"net/http"
	"os"

	"github.com/botopolis/bot"
	"github.com/botopolis/redis"
	"github.com/botopolis/slack"
	"github.com/botopolis/slack/action"
	oslack "github.com/nlopes/slack"
)

func main() {
	adapter := slack.New(os.Getenv("SLACK_TOKEN"))
	r := bot.New(
		adapter,
		redis.New(os.Getenv("REDIS_URL")),
		action.New("/interaction", os.Getenv("SLACK_SIGNATURE")),
	)

	var a action.Plugin
	var slk slack.Adapter
	r.Plugin(&a)
	r.Plugin(&slk)
	a.Add("example", func(cb oslack.AttachmentActionCallback) {
		slk.Client.UpdateMessage(cb.Channel.ID, cb.MessageTs, "updatered")
	})

	r.Hear(bot.Contains("trigger"), func(r bot.Responder) error {
		err := r.Send(bot.Message{
			Params: oslack.PostMessageParameters{
				Attachments: []oslack.Attachment{{
					Text:       "Trigger example",
					CallbackID: "example",
					Actions: []oslack.AttachmentAction{{
						Name:  "check",
						Type:  "button",
						Text:  "Do it",
						Value: "true",
					}, {
						Name:  "check",
						Type:  "button",
						Text:  "Nah",
						Style: "danger",
						Value: "false",
					}},
				}},
			},
		})
		r.Logger.Error("NOPE: %v", err)
		return err
	})

	r.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	}).Methods("GET")

	r.Run()
}
