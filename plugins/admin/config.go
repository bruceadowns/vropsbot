package plugins

import "github.com/bruceadowns/vropsbot/comm"

func init() {
	// bot.Register(&config{})
}

// config is a vropsbot plugin that returns vropsbot configuration
type config struct {
	base
}

// Name returns config name
func (p *config) Name() string {
	return "Config"
}

// Usage returns config usage
func (p *config) Usage() string {
	return "admin config"
}

// Handles returns T/F whether it handles the given directives
func (p *config) Handles(req *comm.Request) (res bool) {
	return
}

// Handle config message
func (p *config) Handle(req *comm.Request) error {
	return nil
}
