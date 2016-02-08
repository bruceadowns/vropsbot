package plugins

import (
	"strings"

	"github.com/bruceadowns/vropsbot/comm"
)

func init() {
	// bot.Register(&report{})
}

// report is a vropsbot plugin that returns
// a daily report about the cat builds
type report struct {
	base
}

// Name returns report name
func (plugin *report) Name() string {
	return "Report"
}

// Usage returns report usage
func (plugin *report) Usage() string {
	return "cat report [branch]"
}

// Handles returns T/F whether it handles the given directives
func (plugin *report) Handles(req *comm.Request) bool {
	if len(req.Arguments) < 2 ||
		len(req.Arguments) > 3 {
		return false
	}

	if !strings.EqualFold(req.Arguments[0], "cat") ||
		!strings.EqualFold(req.Arguments[1], "report") {
		return false
	}

	return true
}

// Handle report message
func (plugin *report) Handle(req *comm.Request) error {
	return nil
}
