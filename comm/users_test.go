package comm

import (
	"os"
	"testing"

	"github.com/nlopes/slack"
)

var utests = []struct {
	ID   string
	Name string
}{
	{
		ID:   "U095NU326",
		Name: "bdowns",
	},
	{
		ID:   "U09KY9MED",
		Name: "ggibson",
	}}

func initUserCache(t *testing.T) {
	a := os.Getenv("SLACK_AUTH_TOKEN")
	if a == "" {
		t.Skip("SLACK_AUTH_TOKEN not provided")
	}
	c := slack.New(a)

	if err := InitUserCache(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetUser(t *testing.T) {
	initUserCache(t)

	for _, v := range utests {
		u, err := GetUser(v.ID)
		if err != nil {
			t.Fatal(err)
		}

		if u != v.Name {
			t.Fatalf("Error for %s. Expected '%s', Actual '%s'\n", v.ID, v.Name, u)
		}
	}
}
