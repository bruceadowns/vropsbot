package plugins

import (
	"strings"

	"github.com/bruceadowns/vropsbot/bot"
	"github.com/bruceadowns/vropsbot/comm"
)

func init() {
	bot.Register(&branches{})
}

// branches is a vropsbot plugin that returns all known cat branches
type branches struct {
	base
}

// Name returns branches name
func (plugin *branches) Name() string {
	return "Branches"
}

// Usage returns branches usage
func (plugin *branches) Usage() string {
	//return "cat branches"
	return ""
}

// Handles returns T/F whether it handles the given directives
func (plugin *branches) Handles(req *comm.Request) bool {
	if len(req.Arguments) != 2 {
		return false
	}

	if !strings.EqualFold(req.Arguments[0], "cat") ||
		!strings.EqualFold(req.Arguments[1], "branches") {
		return false
	}

	return true
}

// Handle branches message
func (plugin *branches) Handle(req *comm.Request) (err error) {
	t := comm.Mention(req.User)

	if len(plugin.config.Branches) == 0 {
		t += "no known branches"
	} else {
		t += "known branches - "
		for k, v := range plugin.config.Branches {
			if k > 0 {
				t += ", "
			}
			t += comm.BuildwebBranchURL(v)
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

	return
}
