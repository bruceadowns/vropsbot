package plugins

import (
	"fmt"
	"strings"

	"github.com/bruceadowns/vropsbot/bot"
	"github.com/bruceadowns/vropsbot/comm"
)

func init() {
	bot.Register(&slas{})
}

// slas is a vropsbot plugin that returns all known cat slas
type slas struct {
	base
}

// Name returns slas name
func (plugin *slas) Name() string {
	return "SLAs"
}

// Usage returns slas usage
func (plugin *slas) Usage() string {
	//return "cat slas"
	return ""
}

// Handles returns T/F whether it handles the given directives
func (plugin *slas) Handles(req *comm.Request) bool {
	if len(req.Arguments) != 2 {
		return false
	}

	if !strings.EqualFold(req.Arguments[0], "cat") ||
		!strings.EqualFold(req.Arguments[1], "slas") {
		return false
	}

	return true
}

// Handle slas message
func (plugin *slas) Handle(req *comm.Request) error {
	t := comm.Mention(req.User)

	if len(plugin.config.SLAs) < 1 {
		t += "no known slas"
	} else {
		t += fmt.Sprintf("known slas - *%s*", strings.Join(plugin.config.SLAs, ", "))
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
