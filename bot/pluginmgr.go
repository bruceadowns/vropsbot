package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bruceadowns/vropsbot/comm"
	"github.com/bruceadowns/vropsbot/common"
)

// PluginManager coalesces all registered plugins
type PluginManager struct {
	pushPlugin []PluginPush
	pullPlugin []PluginPull
}

// AddPush will add a push plugin to the member slice
func (pm *PluginManager) AddPush(p PluginPush) {
	pm.pushPlugin = append(pm.pushPlugin, p)
}

// AddPull will add a pull plugin to the member slice
func (pm *PluginManager) AddPull(p PluginPull) {
	pm.pullPlugin = append(pm.pullPlugin, p)
}

// Prepare calls Prepare on all registered plugins
func (pm *PluginManager) Prepare() error {
	log.Printf("Prepare all pull plugins")
	for _, v := range pm.pullPlugin {
		if err := v.Prepare(); err != nil {
			return err
		}
	}

	log.Printf("Prepare all push plugins")
	for _, v := range pm.pushPlugin {
		if err := v.Prepare(); err != nil {
			return err
		}
	}

	log.Printf("All plugins have been prepared")

	return nil
}

// Start calls Run on all registered push plugins
func (pm *PluginManager) Start() error {
	for _, v := range pm.pushPlugin {
		common.RandomDelay()
		log.Printf("Starting push plugin: %s", v.Name())
		v.Run()
	}
	log.Printf("All push plugins have been started")

	return nil
}

// HandlesHelp checks if help is requested
func (pm *PluginManager) HandlesHelp(req *comm.Request) bool {
	if len(req.Arguments) < 1 {
		return true
	}

	cmd := strings.ToLower(req.Arguments[0])
	cmd = strings.TrimLeft(cmd, "-")
	cmd = strings.TrimLeft(cmd, "—")
	cmd = strings.TrimLeft(cmd, "/")
	switch cmd {
	case "help", "hello", "usage", "h", "?":
		return true
	}

	return false
}

// HandleHelp checks if help is requested
func (pm *PluginManager) HandleHelp(req *comm.Request) error {
	text := comm.Mention(req.User)
	text += "*vR Ops Bot Usage*\n"
	for _, v := range pm.pullPlugin {
		u := v.Usage()
		if len(u) > 0 {
			text += fmt.Sprintf("\t%s\n", v.Usage())
		}
	}

	comm.SB.ChanResponse <- &comm.Response{
		Channel:    req.Channel,
		Text:       text,
		Parameters: comm.DefaultMessageParameters()}

	return nil
}

// HandleUnknown handles unknown directives
func (pm *PluginManager) HandleUnknown(req *comm.Request) error {
	text := comm.Mention(req.User)
	text += fmt.Sprintf("unknown request - *%s*", strings.Join(req.Arguments, " "))

	comm.SB.ChanResponse <- &comm.Response{
		Channel:    req.Channel,
		Text:       text,
		Parameters: comm.DefaultMessageParameters()}

	return nil
}

// HandleError handles plugin handle error
func (pm *PluginManager) HandleError(req *comm.Request, err error) error {
	text := comm.Mention(req.User)
	text += fmt.Sprintf("error occurred - *%s*", err)

	comm.SB.ChanResponse <- &comm.Response{
		Channel:    req.Channel,
		Text:       text,
		Parameters: comm.DefaultMessageParameters()}

	return nil
}

// Dispatch sends arguments to registered pull plugins
func (pm *PluginManager) Dispatch(req *comm.Request) error {
	if pm.HandlesHelp(req) {
		return pm.HandleHelp(req)
	}

	h := 0
	for _, v := range pm.pullPlugin {
		if v.Handles(req) {
			h++
			if err := v.Handle(req); err != nil {
				_ = pm.HandleError(req, err)
				return err
			}
		}
	}

	if h < 1 {
		return pm.HandleUnknown(req)
	}

	return nil
}
