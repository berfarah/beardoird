package uniqlo_test

import (
	"testing"

	"github.com/berfarah/beardroid/plugins/inventory/uniqlo"
	"github.com/stretchr/testify/assert"
)

func TestProductSKU(t *testing.T) {
	expected := "51234COL02SMA001000"
	p := uniqlo.Product{
		ID:    "51234",
		Color: "COL02",
		Size:  "SMA001",
	}

	assert.Equal(t, expected, p.SKU(), "puts the SKU together correctly")
}

func TestProductURL(t *testing.T) {
	expected := "https://www.uniqlo.com/us/en/51234.html?dwvar_51234_color=COL02&dwvar_51234_size=SMA001"
	p := uniqlo.Product{
		ID:    "51234",
		Color: "COL02",
		Size:  "SMA001",
	}

	assert.Equal(t, expected, p.URL(), "puts the URL together correctly")
}

func TestUniqloStockCount(t *testing.T) {
	s := uniqlo.UniqloStock{StrCount: "5"}
	assert.Equal(t, 5, s.Count(), "converts the number")
}

func TestUniqloStockLoad(t *testing.T) {
	var noStock = []byte(`{
	"inStock": false,
	"ats": "0",
	"inStockDate": "",
	"availableForSale": false
}`)

	var inStock = []byte(`{
	"inStock":true,
	"ats":"5",
	"inStockDate":"",
	"availableForSale":false
}`)

	var none uniqlo.UniqloStock
	none.Load(noStock)
	assert.Equal(t, 0, none.Count())
	assert.False(t, none.InStock)

	var some uniqlo.UniqloStock
	some.Load(inStock)
	assert.Equal(t, 5, some.Count())
	assert.True(t, some.InStock)
}
