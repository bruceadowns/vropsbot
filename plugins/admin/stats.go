package plugins

import "github.com/bruceadowns/vropsbot/comm"

func init() {
	// bot.Register(&stats{})
}

// stats is a vropsbot plugin that returns stats and telemetry for vropsbot
type stats struct {
	base
}

// Name returns stats's name
func (p *stats) Name() string {
	return "Stats"
}

// Usage returns stats's usage
func (p *stats) Usage() string {
	return "admin stats"
}

// Handles returns T/F whether it handles the given directives
func (p *stats) Handles(req *comm.Request) (res bool) {
	return
}

// Handle stats message
func (p *stats) Handle(req *comm.Request) error {
	return nil
}
