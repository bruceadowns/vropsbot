package bot

import "log"

// PM is a global PluginManager
var PM PluginManager

// Register plugin with vropsbot
func Register(p Plugin) {
	switch p.(type) {
	case PluginPush:
		log.Printf("Register push plugin %s", p.Name())

		// cast to PluginPush via type assertion
		PM.AddPush(p.(PluginPush))
	case PluginPull:
		log.Printf("Register pull plugin %s", p.Name())

		// cast to PluginPull via type assertion
		PM.AddPull(p.(PluginPull))
	default:
		log.Printf("Unknown plugin: %s", p.Name())
	}
}
