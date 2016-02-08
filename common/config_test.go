package common

import (
	"os"
	"testing"
	"time"
)

func TestStrToSeconds(t *testing.T) {
	var tt = []struct {
		In  string
		Out time.Duration
	}{
		{
			In:  "60",
			Out: 60 * time.Second,
		},
		{
			In:  "60s",
			Out: 60 * time.Second,
		},
		{
			In:  "60m",
			Out: 60 * time.Minute,
		},
		{
			In:  "60h",
			Out: 60 * time.Hour,
		},
	}

	for _, test := range tt {
		d, err := StrToDuration(test.In)
		if err != nil {
			t.Fatal(err)
		}
		if d != test.Out {
			t.Fatalf("In: %s Expected: %v Actual: %d", test.In, test.Out, d)
		}
	}
}

func TestConfig(t *testing.T) {
	n := os.Getenv("VROPSBOT_CONFIG_JSON")
	if n == "" {
		t.Skip("Specify VROPSBOT_CONFIG_JSON")
	}

	c, err := NewConfig()
	if err != nil {
		t.Error(err)
	}

	for _, v := range c.Branches {
		t.Logf("Branch: %s", v)
	}
	t.Logf("Debug: %t", c.Debug)
	t.Logf("SlackbotAuthToken: %s", c.SlackbotAuthToken)
	for _, v := range c.Targets {
		t.Logf("Target: %s", v)
	}
}
