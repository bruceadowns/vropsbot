package plugins

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bruceadowns/vropsbot/bot"
	"github.com/bruceadowns/vropsbot/comm"

	"gopkg.in/jmcvetta/napping.v3"
)

const (
	defaultArea   = "SingleVA_PostCheckin"
	defaultBranch = "main"
	defaultResult = "PASS"
	defaultTarget = "vcopssuitevm"

	unknownStatus = "UNKNOWN"
)

func init() {
	bot.Register(&build{})
}

// build is a vropsbot plugin that returns
// information about the given build
type build struct {
	base
}

// Name returns build name
func (plugin *build) Name() string {
	return "Build"
}

// Usage returns build usage
func (plugin *build) Usage() string {
	return "cat build [branch] [target] [result]"
}

// Handles returns T/F whether it handles the given directives
func (plugin *build) Handles(req *comm.Request) bool {
	if len(req.Arguments) < 2 ||
		len(req.Arguments) > 5 {
		return false
	}

	if !strings.EqualFold(req.Arguments[0], "cat") ||
		!strings.EqualFold(req.Arguments[1], "build") {
		return false
	}

	return true
}

// Handle build message
// CAT URL:
// http://cat.eng.vmware.com/api/v2.0/deliverable/?format=json&limit=1&order_by=-endtime&build__branch=<1>&targets=<2>&result=<3>
func (plugin *build) Handle(req *comm.Request) error {
	url := plugin.config.URLs.BaseCatURL + plugin.config.URLs.BaseBuildURL

	params := DefaultCatParams()

	params["build__branch"] = defaultBranch
	if len(req.Arguments) > 2 {
		params["build__branch"] = req.Arguments[2]
	}

	params["targets"] = defaultTarget
	if len(req.Arguments) > 3 {
		params["targets"] = req.Arguments[3]
	}

	if len(req.Arguments) > 4 {
		params["result"] = req.Arguments[4]
	}

	p := params.AsUrlValues()
	s := napping.Session{}
	j := BuildJSON{}
	r, err := s.Get(url, &p, &j, nil)
	if err != nil {
		return err
	}
	if r.Status() != 200 {
		return fmt.Errorf("expect response status code 200. Actual %d %s", r.Status(), params)
	}
	if len(j.Objects) == 0 {
		return fmt.Errorf("build %s/%s not found", params["build__branch"], params["targets"])
	}
	if len(j.Objects) > 1 {
		return fmt.Errorf("expected one object. Actual %d %s", len(j.Objects), params)
	}

	t := comm.Mention(req.User)
	t += fmt.Sprintf("latest *%s/%s* ", j.Objects[0].Build.Branch, j.Objects[0].Targets)
	t += fmt.Sprintf("status is *%s* ", j.Objects[0].Result)
	t += fmt.Sprintf("(%s)", comm.BuildwebSbURL(j.Objects[0].SbBuildID, strconv.Itoa(j.Objects[0].SbBuildID)))

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
