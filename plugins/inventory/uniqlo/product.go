package uniqlo

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	XXS size = "SMA001"
	XS  size = "SMA002"
	S   size = "SMA003"
	M   size = "SMA004"
	L   size = "SMA005"
	XL  size = "SMA006"
	XXL size = "SMA007"
)

type size string

type UniqloStock struct {
	InStock  bool   `json:"inStock"`
	StrCount string `json:"ats"`
}

func (s UniqloStock) Count() int {
	i, _ := strconv.Atoi(s.StrCount)
	return i
}

type Product struct {
	ID    string
	Color string
	Size  size
}

func (p Product) SKU() string {
	return fmt.Sprintf("%s%s%s000", p.ID, p.Color, p.Size)
}

const productURL = "https://www.uniqlo.com/us/en/%s.html?dwvar_%s_color=%s&dwvar_%s_size=%s"

func (p Product) URL() string {
	return fmt.Sprintf(
		productURL,
		p.ID,
		p.ID, p.Color,
		p.ID, p.Size,
	)
}

func (p Product) Available() bool {
	b, err := (&Request{SKU: p.SKU()}).Do()
	if err != nil {
		return false
	}
	var s UniqloStock
	if err := json.Unmarshal(b, &s); err != nil {
		return false
	}

	return s.InStock
}
