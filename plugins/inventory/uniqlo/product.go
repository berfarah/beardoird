package uniqlo

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type UniqloStock struct {
	InStock  bool   `json:"inStock"`
	StrCount string `json:"ats"`
}

func (s *UniqloStock) Count() int {
	i, _ := strconv.Atoi(s.StrCount)
	return i
}

func (s *UniqloStock) Load(b []byte) { json.Unmarshal(b, s) }

type Product struct {
	ID    string
	Color string
	Size  string
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
	s.Load(b)

	return s.InStock
}
