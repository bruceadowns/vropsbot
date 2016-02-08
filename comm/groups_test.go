package comm

import (
	"os"
	"testing"

	"github.com/nlopes/slack"
)

var gtests = []struct {
	id   string
	name string
}{
	{
		id:   "G0BVAS72M",
		name: "vrops-bad",
	}}

func initGroupCache(t *testing.T) {
	a := os.Getenv("SLACK_AUTH_TOKEN")
	if a == "" {
		t.Skip("SLACK_AUTH_TOKEN not provided")
	}
	c := slack.New(a)
	if err := InitGroupCache(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetGroup(t *testing.T) {
	initGroupCache(t)

	for _, v := range gtests {
		n, err := GetGroup(v.id)
		if err != nil {
			t.Error(err)
		}
		if n != v.name {
			t.Errorf("%s Expected %s Actual %s", v.id, v.name, n)
		}
	}
}

func TestGetGroupByName(t *testing.T) {
	initGroupCache(t)

	for _, v := range gtests {
		id, err := GetGroupByName(v.name)
		if err != nil {
			t.Error(err)
		}
		if id != v.id {
			t.Errorf("%s Expected %s Actual %s", v.name, v.id, id)
		}
	}
}
