package uniqlo

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/nlopes/slack"
)

type color struct{ Code, Name string }

type colorCheck struct {
	ID string
}

func (c colorCheck) Request() (io.Reader, error) {
	res, err := http.Get(Product{ID: c.ID}.URL())
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("Non-200 status code")
	}

	return res.Body, nil
}

func (c colorCheck) Parse(read io.Reader) (colors []color) {
	b, _ := ioutil.ReadAll(read)
	expr := fmt.Sprintf(`img src="https://uniqlo\.scene7\.com/is/image/UNIQLO/goods_(\d*)_%s_chip\?\$small_ec\$" alt="([\w\s]*)"`, c.ID)
	r := regexp.MustCompile(expr)

	matches := r.FindAllStringSubmatch(string(b), -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		colors = append(colors, color{Code: "COL" + match[1], Name: match[2]})
	}

	return
}

func getColors(id string) (colors []color, err error) {
	c := colorCheck{ID: id}
	b, err := c.Request()
	if err != nil {
		return colors, err
	}
	return c.Parse(b), nil
}

func colorOptions(colors []color) (opts []slack.AttachmentActionOption) {
	for _, color := range colors {
		opts = append(opts, slack.AttachmentActionOption{Text: color.Name, Value: color.Code})
	}
	return opts
}

func slackOptions(colors []color) slack.Attachment {
	return slack.Attachment{
		Text:       "Choose your size and color:",
		Fallback:   ":(",
		CallbackID: "item_selection",
		Actions: []slack.AttachmentAction{
			{
				Name:    "size",
				Text:    "Pick your size",
				Type:    "select",
				Options: sizeOptions,
			},
			{
				Name:    "color",
				Text:    "Pick your color",
				Type:    "select",
				Options: colorOptions(colors),
			},
		},
	}
}
