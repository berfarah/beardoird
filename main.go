package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/berfarah/beardroid/ludlow"
	"github.com/berfarah/gobot"
)

func main() {
	var secret string
	if contents, err := ioutil.ReadFile(".secret"); err != nil {
		panic(err)
	} else {
		secret = strings.TrimSpace(string(contents))
	}

	// logger := log.New(os.Stdout, "beardroid: ", log.Lshortfile|log.LstdFlags)
	r := gobot.New(secret)
	r.Load(ludlow.Plugin)
	r.Hear(&gobot.Hook{
		Name:    "yoyo",
		Matcher: gobot.MatchText("heyy"),
		Func: func(r *gobot.Responder) error {
			fmt.Println("running thing")
			r.Send(gobot.Message{Text: "ayy"})
			return nil
		},
	})
	r.Connect()
}
