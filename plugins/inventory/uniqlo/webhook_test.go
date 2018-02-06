package uniqlo

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/nlopes/slack"
)

func TestParseJSON(t *testing.T) {
	testData, err := os.Open("testdata/request.json")
	if err != nil {
		t.Errorf("Couldn't open testdata: %v", err)
	}
	b, err := ioutil.ReadAll(testData)
	if err != nil {
		t.Errorf("%v", err)
	}

	var cb slack.AttachmentActionCallback
	json.Unmarshal(b, &cb)

	t.Errorf("%v", cb)
}
