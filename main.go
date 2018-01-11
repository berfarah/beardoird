package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/berfarah/beardroid/brain"
	slackerino "github.com/berfarah/beardroid/slack"
)

func main() {
	var secret string
	if contents, err := ioutil.ReadFile(".secret"); err != nil {
		panic(err)
	} else {
		secret = strings.TrimSpace(string(contents))
	}

	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	adapter := slackerino.New(secret, logger)
	brain.Write(nil)
	adapter.Listen()
}
