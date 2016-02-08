package plugins

import (
	"fmt"
	"time"

	"github.com/bruceadowns/vropsbot/bot"
	"github.com/bruceadowns/vropsbot/comm"
)

func init() {
	bot.Register(&uptime{})
}

// uptime is a vropsbot plugin that returns vropsbot uptime
type uptime struct {
	base
	start time.Time
}

// Name returns uptime's name
func (plugin *uptime) Name() string {
	return "Uptime"
}

// Usage returns uptime's usage
func (plugin *uptime) Usage() string {
	return "admin uptime"
}

// Prepare initializes uptime's context
func (plugin *uptime) Prepare() error {
	plugin.start = time.Now()

	return nil
}

// Handles returns T/F whether it handles the given directives
func (plugin *uptime) Handles(req *comm.Request) bool {
	if len(req.Arguments) == 1 &&
		req.Arguments[0] == "admin" {
		return true
	}

	if len(req.Arguments) == 2 &&
		req.Arguments[0] == "admin" &&
		req.Arguments[1] == "uptime" {
		return true
	}

	return false
}

// Handle uptime message
func (plugin *uptime) Handle(req *comm.Request) error {
	now := time.Now()

	t := comm.Mention(req.User)
	t += fmt.Sprintf(
		"start time: *%s* current time: *%s* uptime: *%s*",
		plugin.start.Format(time.RFC822),
		now.Format(time.RFC822),
		now.Sub(plugin.start))

	comm.SB.ChanResponse <- &comm.Response{
		Channel:    req.Channel,
		Text:       t,
		Parameters: comm.DefaultMessageParameters()}

	return nil
}
