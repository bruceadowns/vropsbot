package comm

import (
	"os"
	"testing"

	"github.com/nlopes/slack"
)

var ctests = []struct {
	ID   string
	Name string
}{
	{
		ID:   "C092SNWRE",
		Name: "vrops-general",
	},
	{
		ID:   "C098GP2R3",
		Name: "vrops-install",
	},
}

func initChannelCache(t *testing.T) {
	a := os.Getenv("SLACK_AUTH_TOKEN")
	if a == "" {
		t.Skip("SLACK_AUTH_TOKEN not provided")
	}
	c := slack.New(a)
	if err := InitChannelCache(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetChannel(t *testing.T) {
	initChannelCache(t)

	for _, v := range ctests {
		n, err := GetChannel(v.ID)
		if err != nil {
			t.Error(err)
		}
		if n != v.Name {
			t.Errorf("%s Expected %s Actual %s", v.ID, v.Name, n)
		}
	}
}

func TestGetChannelByName(t *testing.T) {
	initChannelCache(t)

	for _, v := range ctests {
		id, err := GetChannelByName(v.Name)
		if err != nil {
			t.Error(err)
		}
		if id != v.ID {
			t.Errorf("%s Expected %s Actual %s", v.Name, v.ID, id)
		}
	}
}
