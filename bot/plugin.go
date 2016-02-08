package bot

import (
	"github.com/bruceadowns/vropsbot/comm"
)

// Plugin interface to be implemented via pull plugins
type Plugin interface {
	Name() string
	Prepare() error
}

// PluginPush interface to be implemented via push plugins
type PluginPush interface {
	Plugin
	Run() error
}

// PluginPull interface to be implemented via pull plugins
type PluginPull interface {
	Plugin
	Usage() string
	Handles(req *comm.Request) bool
	Handle(req *comm.Request) error
}
