package plugins

import (
	"strings"

	"github.com/bruceadowns/vropsbot/bot"
	"github.com/bruceadowns/vropsbot/comm"
)

func init() {
	bot.Register(&areas{})
}

// areas is a vropsbot plugin that returns all known cat areas
type areas struct {
	base
}

// Name returns areas name
func (plugin *areas) Name() string {
	return "Areas"
}

// Usage returns areas usage
func (plugin *areas) Usage() string {
	//return "cat areas"
	return ""
}

// Handles returns T/F whether it handles the given directives
func (plugin *areas) Handles(req *comm.Request) bool {
	if len(req.Arguments) != 2 {
		return false
	}

	if !strings.EqualFold(req.Arguments[0], "cat") ||
		!strings.EqualFold(req.Arguments[1], "areas") {
		return false
	}

	return true
}

// Handle areas message
func (plugin *areas) Handle(req *comm.Request) error {
	t := comm.Mention(req.User)

	if len(plugin.config.Areas) < 1 {
		t += "no known areas"
	} else {
		t += "known areas - "
		for k, v := range plugin.config.Areas {
			if k > 0 {
				t += ", "
			}
			t += comm.CatAreaURL(v)
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
