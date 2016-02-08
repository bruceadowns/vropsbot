package plugins

import (
	"fmt"
	"strings"

	"github.com/bruceadowns/vropsbot/bot"
	"github.com/bruceadowns/vropsbot/comm"
)

func init() {
	bot.Register(&welcome{})
}

// welcome is a vropsbot plugin that returns vropsbot configuration
type welcome struct {
	base
}

var (
	responseMap = map[string]string{
		"thank":    "you are welcome",
		"thanks":   "you're welcome",
		"thankyou": "yourewelcome",
		"ty":       "yw",
		"gracias":  "de nada",
		"merci":    "de rein",
		"danke":    "gern geschehen",
	}

	replacer = strings.NewReplacer(" ", "", ".", "", "!", "")
)

// Name returns welcome name
func (plugin *welcome) Name() string {
	return "Welcome"
}

// Usage returns welcome usage
func (plugin *welcome) Usage() string {
	//return "thank you"
	return ""
}

// Handles returns T/F whether it handles the given directives
func (plugin *welcome) Handles(req *comm.Request) bool {
	if len(req.Arguments) == 1 {
		m := strings.ToLower(req.Arguments[0])
		m = replacer.Replace(m)

		if _, ok := responseMap[m]; ok {
			return true
		}
	}

	if len(req.Arguments) == 2 &&
		strings.EqualFold(req.Arguments[0], "thank") &&
		strings.EqualFold(req.Arguments[1], "you") {
		return true
	}

	return false
}

// Handle welcome message
func (plugin *welcome) Handle(req *comm.Request) error {
	t := comm.Mention(req.User)
	m := strings.ToLower(req.Arguments[0])
	m = replacer.Replace(m)
	r, _ := responseMap[m]
	t += fmt.Sprintf("%s", r)

	comm.SB.ChanResponse <- &comm.Response{
		Channel:    req.Channel,
		Text:       t,
		Parameters: comm.DefaultMessageParameters()}

	return nil
}
