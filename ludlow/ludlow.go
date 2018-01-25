package ludlow

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/berfarah/gobot"
	"github.com/nlopes/slack"
)

type main struct{}

// Plugin is the exportable plugin
var Plugin = main{}

const productCode = "H2327"
const url = "https://www.jcrew.com/p/mens_category/outerwear/topcoats/ludlow-topcoat-in-italian-woolcashmere-with-thinsulate/H2327"

type Product struct {
	Color    string
	SKU      string
	Quantity int
}

var productBySku = map[string]*Product{
	"99104786265": &Product{"Navy", "99104786265", 0},
	// "99104786301": &Product{"Charcoal", "99104786301", 0},
}
var client = &http.Client{}

func (m main) Load(r *gobot.Robot) {
	ticker := time.NewTicker(1 * time.Hour)
	checkStock(r)
	for range ticker.C {
		checkStock(r)
	}
}

func checkStock(r *gobot.Robot) {
	decodeRequest()
	if msg, ok := stockMessage(); ok {
		r.Logger.Debug("Ludlow: Items in stock")
		r.Chat.Send(gobot.Message{
			Room:   "@bernardo",
			Text:   "",
			Params: msg,
		})
		return
	}
	r.Logger.Debug("Ludlow: No items in stock")
}

func stockMessage() (msg slack.PostMessageParameters, ok bool) {
	msg.Attachments = []slack.Attachment{{
		Title:     "Ludlow is in stock!",
		TitleLink: url,
		Fallback:  "Ludlow is in stock!",
		Color:     "good",
		Fields:    make([]slack.AttachmentField, len(productBySku)),
	}}

	var i int
	for _, p := range productBySku {
		if p.Quantity > 0 {
			ok = true
		}
		msg.Attachments[0].Fields[i] = slack.AttachmentField{
			Title: p.Color,
			Value: strconv.Itoa(p.Quantity),
			Short: true,
		}
		i++
	}
	return msg, ok
}

func decodeRequest() {
	r, err := request()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	d := json.NewDecoder(r)
	var i interface{}
	if err := d.Decode(&i); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	top := i.(map[string]interface{})
	inventory := top["inventory"].(map[string]interface{})
	for sku, stock := range inventory {
		if p, ok := productBySku[sku]; ok {
			s := stock.(map[string]interface{})
			if q, ok := s["quantity"].(float64); ok {
				p.Quantity = int(q)
			}
		}
	}
}

func request() (io.Reader, error) {
	var b io.Reader
	req, err := http.NewRequest("GET", "https://www.jcrew.com/data/v1/US/products/inventory/"+productCode, nil)
	if err != nil {
		return b, err
	}
	addHeaders(req)

	res, err := client.Do(req)
	if err != nil {
		return b, err
	}
	if res.StatusCode != 200 {
		return b, fmt.Errorf("API returned a %s", res.Status)
	}

	return gzip.NewReader(res.Body)
}

func addHeaders(req *http.Request) {
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:57.0) Gecko/20100101 Firefox/57.0")
	req.Header.Add("Referer", "https://www.jcrew.com/p/mens_c%E2%80%A62327?color_name=hthr-charcoal")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Host", "www.jcrew.com")
	req.Header.Add("origin", "https://www.jcrew.com")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Pragma", "no-cache")
}
