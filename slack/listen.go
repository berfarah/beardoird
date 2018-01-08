package slack

import (
	"fmt"
	"log"

	"github.com/nlopes/slack"
)

type Adapter struct {
	api    *slack.Client
	rtm    *slack.RTM
	logger *log.Logger
}

func New(secret string, logger *log.Logger) *Adapter {
	slack.SetLogger(logger)
	return &Adapter{
		api:    slack.New(secret),
		logger: logger,
	}
}

func (a *Adapter) Listen() {
	a.rtm = a.api.NewRTM()
	go a.rtm.ManageConnection()

	cache := NewCache(a.api)
	if err := cache.Populate(); err != nil {
		a.logger.Fatalln(err)
	}

	for msg := range a.rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:

		case *slack.ConnectedEvent:
			params := slack.PostMessageParameters{}
			params.Attachments = []slack.Attachment{
				slack.Attachment{
					Color: "danger",
					Fields: []slack.AttachmentField{
						slack.AttachmentField{
							Title: "foo",
							Value: "bar",
						},
					},
				},
			}

			a.api.PostMessage(cache.DMs["bernardo"], "ayy", params)

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

func (a *Adapter) Disconnect() error {
	return a.rtm.Disconnect()
}
