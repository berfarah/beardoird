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
	// var self *slack.UserInfo
	a.rtm = a.api.NewRTM()
	go a.rtm.ManageConnection()

	cache := newCache(a.api)

	for msg := range a.rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			fmt.Println("Hello")

		case *slack.ConnectedEvent:
			fmt.Printf("Connected: %v\n", ev.ConnectionCount)
			cache.Populate(
				ev.Info.Users,
				ev.Info.Channels,
				ev.Info.IMs,
			)
		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

			if ev.Text != "heyy" {
				continue
			}

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

			dm := DM{Username: "bernardo"}
			if err := cache.brain.Read(&dm); err != nil {
				a.logger.Fatal(err)
				return
			}
			a.api.PostMessage(dm.Conversation, "ayy", params)
			a.rtm.SendMessage(a.rtm.NewOutgoingMessage("ayy", ev.Channel))

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
			fmt.Printf("Unexpected: %v\n", msg.Data)

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

func (a *Adapter) Disconnect() error {
	return a.rtm.Disconnect()
}
