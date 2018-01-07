package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

func main() {
	secret, err := ioutil.ReadFile(".secret")
	if err != nil {
		panic(err)
	}

	api := slack.New(strings.TrimSpace(string(secret)))
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	cache := NewCache(api)
	if err := cache.Populate(); err != nil {
		logger.Fatalln(err)
	}

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)

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

			api.PostMessage(cache.DMs["bernardo"], "ayy", params)

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
