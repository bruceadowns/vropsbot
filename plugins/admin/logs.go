package plugins

import "github.com/bruceadowns/vropsbot/comm"

func init() {
	// bot.Register(&logs{})
}

// logs is a vropsbot plugin that returns vropsbot log files
type logs struct {
	base
}

// Name returns logs name
func (p *logs) Name() string {
	return "Logs"
}

// Usage returns logs usage
func (p *logs) Usage() string {
	return "admin logs"
}

// Handles returns T/F whether it handles the given directives
func (p *logs) Handles(req *comm.Request) (res bool) {
	return
}

// Handle logs message
func (p *logs) Handle(req *comm.Request) error {
	return nil
}
