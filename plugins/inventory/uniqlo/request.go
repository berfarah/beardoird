package uniqlo

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Req *http.Request
	SKU string
}

const inventoryURL = "https://www.uniqlo.com/on/demandware.store/Sites-UniqloUS-Site/default/Product-GetAvailability?pid=%s&Quantity=1"

func (r *Request) generate() {
	req, _ := http.NewRequest("GET", fmt.Sprintf(inventoryURL, r.SKU), nil)
	req.Header.Set("Host", "www.uniqlo.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:58.0) Gecko/20100101 Firefox/58.0")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")

	r.Req = req
}

func (r *Request) Do() ([]byte, error) {
	r.generate()
	res, err := http.DefaultClient.Do(r.Req)
	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode != 200 {
		return []byte{}, fmt.Errorf("API returned a %s", res.Status)
	}

	zip, err := gzip.NewReader(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return ioutil.ReadAll(zip)
}
