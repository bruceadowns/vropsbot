package plugins

import "github.com/bruceadowns/vropsbot/comm"

func init() {
	// bot.Register(&reset{})
}

// reset is a vropsbot plugin that resets the vropsbot
type reset struct {
	base
}

// Name returns reset's name
func (p *reset) Name() string {
	return "Reset"
}

// Usage returns reset's usage
func (p *reset) Usage() string {
	return "admin reset"
}

// Handles returns T/F whether it handles the given directives
func (p *reset) Handles(req *comm.Request) (res bool) {
	return
}

// Handle reset message
func (p *reset) Handle(req *comm.Request) error {
	return nil
}
