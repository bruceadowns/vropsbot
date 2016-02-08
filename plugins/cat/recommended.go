package plugins

import (
	"fmt"
	"strings"

	"github.com/bruceadowns/vropsbot/bot"
	"github.com/bruceadowns/vropsbot/comm"

	"gopkg.in/jmcvetta/napping.v3"
)

func init() {
	bot.Register(&recommended{})
}

// recommended is a vropsbot plugin that returns
// information about the given recommended build
type recommended struct {
	base
}

// Name returns recommended name
func (plugin *recommended) Name() string {
	return "Recommended"
}

// Usage returns recommended usage
func (plugin *recommended) Usage() string {
	return "cat recommended [branch]"
}

// Handles returns T/F whether it handles the given directives
func (plugin *recommended) Handles(req *comm.Request) bool {
	if len(req.Arguments) < 2 ||
		len(req.Arguments) > 3 {
		return false
	}

	if !strings.EqualFold(req.Arguments[0], "cat") ||
		!strings.EqualFold(req.Arguments[1], "recommended") {
		return false
	}

	return true
}

// Handle recommended message
// CAT URL:
// http://cat.eng.vmware.com/api/v2.0/recommendation/?format=json&limit=1&order_by=-updated&branch__name=main&sla__name=VA_Bats&site=mbu
func (plugin *recommended) Handle(req *comm.Request) error {
	branch := defaultBranch
	if len(req.Arguments) > 2 {
		branch = req.Arguments[2]
	}

	jRecommended := RecommendedJSON{}
	{
		url := plugin.config.URLs.BaseCatURL + plugin.config.URLs.BaseRecommendedURL
		params := DefaultRecommendedParams()
		params["branch__name"] = branch

		p := params.AsUrlValues()
		s := napping.Session{}
		r, err := s.Get(url, &p, &jRecommended, nil)
		if err != nil {
			return err
		}
		if r.Status() != 200 {
			return fmt.Errorf("expect response status code 200. Actual %d %s", r.Status(), params)
		}
		if len(jRecommended.Objects) == 0 {
			return fmt.Errorf("recommended change list for %s not found", params["branch__name"])
		}
		if len(jRecommended.Objects) > 1 {
			return fmt.Errorf("expected one object. Actual %d %s", len(jRecommended.Objects), params)
		}
	}

	// CAT URL:
	// http://cat.eng.vmware.com/api/v2.0/build/<id>/?format=json
	jCurrBuild := CurrBuildJSON{}
	{
		url := plugin.config.URLs.BaseCatURL + jRecommended.Objects[0].CurrBuild
		params := DefaultCurrBuildParams()

		p := params.AsUrlValues()
		s := napping.Session{}
		r, err := s.Get(url, &p, &jCurrBuild, nil)
		if err != nil {
			return err
		}
		if r.Status() != 200 {
			return fmt.Errorf("expect response status code 200. Actual %d\n", r.Status())
		}
	}

	t := comm.Mention(req.User)
	t += fmt.Sprintf("recommended changeset for *%s* ", branch)
	t += fmt.Sprintf("is *%s*", comm.P4WebBranchURL(branch, jCurrBuild.Changeset))

	iop := &comm.Interop{
		Request: req,
		Response: &comm.Response{
			Channel:    req.Channel,
			Text:       t,
			Parameters: comm.DefaultMessageParameters()}}
	comm.SB.ChanResponse <- iop.Response
	comm.SB.ChanPersist <- iop

	return nil
}
