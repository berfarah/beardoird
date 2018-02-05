package uniqlo

import "github.com/nlopes/slack"

var sizes = map[string]string{
	"XXS": "SMA001",
	"XS":  "SMA002",
	"S":   "SMA003",
	"M":   "SMA004",
	"L":   "SMA005",
	"XL":  "SMA006",
	"XXL": "SMA007",
}

var sizeOptions = []slack.AttachmentActionOption{
	{Text: "XXS", Value: "XXS"},
	{Text: "XS", Value: "XS"},
	{Text: "S", Value: "S"},
	{Text: "M", Value: "M"},
	{Text: "L", Value: "L"},
	{Text: "XL", Value: "XL"},
	{Text: "XXL", Value: "XXL"},
}

const (
	XXS = sizes["XXS"]
	XS  = sizes["XS"]
	S   = sizes["S"]
	M   = sizes["M"]
	L   = sizes["L"]
	XL  = sizes["XL"]
	XXL = sizes["XXL"]
)
