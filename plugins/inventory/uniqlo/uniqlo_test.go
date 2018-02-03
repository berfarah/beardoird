package uniqlo_test

import (
	"testing"

	"github.com/berfarah/beardroid/plugins/uniqlo/uniqlo"
	"github.com/stretchr/testify/assert"
)

var uniqloNone = []byte(`{
	"inStock": false,
	"ats": "0",
	"inStockDate": "",
	"availableForSale": false
}`)

var uniqloInStock = []byte(`{
	"inStock":true,
	"ats":"5",
	"inStockDate":"",
	"availableForSale":false
}`)

func TestUniqloAvailable(t *testing.T) {
	u := uniqlo.Uniqlo{}

	r, _ := u.Available(uniqloNone)
	assert.Equal(t, 0, r.Count())
	assert.False(t, r.InStock)

	r, _ = u.Available(uniqloInStock)
	assert.Equal(t, 5, r.Count())
	assert.True(t, r.InStock)
}
