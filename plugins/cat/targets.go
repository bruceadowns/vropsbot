package plugins

import (
	"strings"

	"github.com/bruceadowns/vropsbot/bot"
	"github.com/bruceadowns/vropsbot/comm"
)

func init() {
	bot.Register(&targets{})
}

// targets is a vropsbot plugin that returns all known cat targets
type targets struct {
	base
}

// Name returns targets name
func (plugin *targets) Name() string {
	return "Targets"
}

// Usage returns targets usage
func (plugin *targets) Usage() string {
	//return "cat targets"
	return ""
}

// Handles returns T/F whether it handles the given directives
func (plugin *targets) Handles(req *comm.Request) bool {
	if len(req.Arguments) != 2 {
		return false
	}

	if !strings.EqualFold(req.Arguments[0], "cat") ||
		!strings.EqualFold(req.Arguments[1], "targets") {
		return false
	}

	return true
}

// Handle targets message
func (plugin *targets) Handle(req *comm.Request) error {
	t := comm.Mention(req.User)

	if len(plugin.config.Targets) < 1 {
		t += "no known targets"
	} else {
		t += "known targets - "
		for k, v := range plugin.config.Targets {
			if k > 0 {
				t += ", "
			}
			t += comm.BuildwebTargetURL(v)
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
