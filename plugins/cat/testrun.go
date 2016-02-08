package plugins

import (
	"fmt"
	"strings"

	"github.com/bruceadowns/vropsbot/bot"
	"github.com/bruceadowns/vropsbot/comm"

	"gopkg.in/jmcvetta/napping.v3"
)

func init() {
	bot.Register(&testrun{})
}

// testrun is a vropsbot plugin that returns
// information about the given testrun
type testrun struct {
	base
}

// Name returns testrun name
func (plugin *testrun) Name() string {
	return "Testrun"
}

// Usage returns testrun usage
func (plugin *testrun) Usage() string {
	return "cat testrun [branch] [area] [result]"
}

// Handles returns T/F whether it handles the given directives
func (plugin *testrun) Handles(req *comm.Request) bool {
	if len(req.Arguments) < 2 || len(req.Arguments) > 5 {
		return false
	}

	if !strings.EqualFold(req.Arguments[0], "cat") ||
		!strings.EqualFold(req.Arguments[1], "testrun") {
		return false
	}

	return true
}

// Handle testrun message
// CAT URL:
// http://cat.eng.vmware.com/api/v2.0/testrun/?format=json&limit=1&order_by=-endtime&deliverables__build__branch=main&area__name=SingleVA_PostCheckin&result=PASS
func (plugin *testrun) Handle(req *comm.Request) error {
	url := plugin.config.URLs.BaseCatURL + plugin.config.URLs.BaseTestrunURL

	params := DefaultCatParams()

	params["deliverables__build__branch"] = defaultBranch
	if len(req.Arguments) > 2 {
		params["deliverables__build__branch"] = req.Arguments[2]
	}

	params["area__name"] = defaultArea
	if len(req.Arguments) > 3 {
		params["area__name"] = req.Arguments[3]
	}

	if len(req.Arguments) > 4 {
		params["result"] = req.Arguments[4]
	}

	p := params.AsUrlValues()
	s := napping.Session{}
	j := TestrunJSON{}
	r, err := s.Get(url, &p, &j, nil)
	if err != nil {
		return err
	}
	if r.Status() != 200 {
		return fmt.Errorf("expect response status code 200. Actual %d %s", r.Status(), params)
	}
	if len(j.Objects) == 0 {
		return fmt.Errorf("testrun %s/%s not found", params["deliverables__build__branch"], params["area__name"])
	}
	if len(j.Objects) != 1 {
		return fmt.Errorf("expected one object. Actual %d %s", len(j.Objects), params)
	}

	t := comm.Mention(req.User)
	t += fmt.Sprintf("testrun for %s ",
		comm.CatTestrunURL(j.Objects[0].ID,
			fmt.Sprintf("%s/%s", params["deliverables__build__branch"], params["area__name"])))
	t += fmt.Sprintf("is *%s* ", j.Objects[0].Result)
	t += fmt.Sprintf("(%s)", comm.CatTestrunResultsURL(j.Objects[0].ResultsDir))

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
