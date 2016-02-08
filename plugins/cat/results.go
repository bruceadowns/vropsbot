package plugins

import (
	"strings"

	"github.com/bruceadowns/vropsbot/bot"
	"github.com/bruceadowns/vropsbot/comm"
)

func init() {
	bot.Register(&results{})
}

// results is a vropsbot plugin that returns all known cat results
type results struct {
	base
}

// Name returns results name
func (plugin *results) Name() string {
	return "Results"
}

// Usage returns results usage
func (plugin *results) Usage() string {
	//return "cat results"
	return ""
}

// Handles returns T/F whether it handles the given directives
func (plugin *results) Handles(req *comm.Request) bool {
	if len(req.Arguments) != 2 {
		return false
	}

	if !strings.EqualFold(req.Arguments[0], "cat") ||
		!strings.EqualFold(req.Arguments[1], "results") {
		return false
	}

	return true
}

// Handle results message
func (plugin *results) Handle(req *comm.Request) error {
	t := comm.Mention(req.User)

	if len(plugin.config.Results) < 1 {
		t += "no known results"
	} else {
		t += "known results - "
		for k, v := range plugin.config.Results {
			if k == 0 {
				t += "*"
			} else {
				t += ", *"
			}
			t += v
			t += "*"
		}
	}

	iop := &comm.Interop{
		Request: req,
		Response: &comm.Response{
			Channel:    req.Channel,
			Text:       t,
			Parameters: comm.DefaultMessageParameters()}}
	comm.SB.ChanResponse <- iop.Response
	comm.SB.ChanPersist <- iop

	return nil
}
