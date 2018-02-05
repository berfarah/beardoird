package uniqlo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseColors(t *testing.T) {
	assert := assert.New(t)
	testData, err := os.Open("testdata/site.html")
	if err != nil {
		t.Errorf("Couldn't open testdata: %v", err)
	}
	c := colorCheck{ID: "401925"}
	colors := c.Parse(testData)

	assert.Equal("COL01", colors[0].Code)
	assert.Equal("OFF WHITE", colors[0].Name)

	assert.Equal("COL09", colors[1].Code)
	assert.Equal("BLACK", colors[1].Name)

	assert.Equal("COL17", colors[2].Code)
	assert.Equal("RED", colors[2].Name)

	assert.Equal("COL69", colors[3].Code)
	assert.Equal("NAVY", colors[3].Name)
}
