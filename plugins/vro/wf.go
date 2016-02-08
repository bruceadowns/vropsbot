package plugins

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bruceadowns/vropsbot/bot"
	"github.com/bruceadowns/vropsbot/comm"

	"gopkg.in/jmcvetta/napping.v3"
)

func init() {
	bot.Register(&wf{})
}

type wf struct {
	base
}

// Name returns the workflow plugin's name
func (plugin *wf) Name() string {
	return "Workflow"
}

// Usage returns the workflow plugin's usage
func (plugin *wf) Usage() string {
	return "vro workflow build [type] [name]"
}

// Handles returns whether this plugin will handle the given parameters
func (plugin *wf) Handles(req *comm.Request) bool {
	if len(req.Arguments) < 3 ||
		len(req.Arguments) > 5 {
		return false
	}

	if !strings.EqualFold(req.Arguments[0], "vro") {
		return false
	}

	if !strings.EqualFold(req.Arguments[1], "workflow") &&
		!strings.EqualFold(req.Arguments[1], "wf") {
		return false
	}

	if _, err := strconv.Atoi(req.Arguments[2]); err != nil {
		return false
	}

	if len(req.Arguments) > 3 &&
		!strings.EqualFold(req.Arguments[3], "ob") &&
		!strings.EqualFold(req.Arguments[3], "sb") {
		return false
	}

	return true
}

// Handle manages the posting of a new vro workflow
// GET https://10.25.37.28:8281/vco/api/workflows/
// POST https://10.25.37.28:8281/vco/api/workflows/<id>/executions
func (plugin *wf) Handle(req *comm.Request) error {
	vroURL := plugin.config.Vro.BaseURL + plugin.config.Vro.WfURL
	vroUser := plugin.config.Vro.Username
	vroPass := plugin.config.Vro.Password
	vroInsecure := plugin.config.Vro.Insecure
	debug := plugin.config.Debug

	build := req.Arguments[2]

	buildType := "sb"
	if len(req.Arguments) > 3 {
		buildType = req.Arguments[3]
	}

	wfName := "dev_checkin_tests"
	if len(req.Arguments) > 4 {
		wfName = req.Arguments[4]
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: vroInsecure},
	}
	s := napping.Session{
		Client:   &http.Client{Transport: tr},
		Log:      debug,
		Userinfo: url.UserPassword(vroUser, vroPass)}

	var href string
	{
		np := napping.Params{
			"conditions": fmt.Sprintf("name=%s", wfName)}
		p := np.AsUrlValues()
		j := WorkflowsJSON{}
		r, err := s.Get(vroURL, &p, &j, nil)
		if err != nil {
			return err
		}
		if r.Status() != 200 {
			return fmt.Errorf("error querying vRO for '%s' [%d]", wfName, r.Status())
		}
		if len(j.Link) == 0 {
			return fmt.Errorf("vRO workflow '%s' not found", wfName)
		}
		if len(j.Link) != 1 {
			return fmt.Errorf("expect one vRO %s workflow. Actual %d", wfName, len(j.Link))
		}
		if len(j.Link[0].Href) == 0 {
			return fmt.Errorf("href is empty for '%s' vRO workflow", wfName)
		}
		href = j.Link[0].Href
	}

	{
		slackUser, _ := comm.GetUser(req.User)
		slackEmail := fmt.Sprintf("%s@vmware.com", slackUser)

		j := WorkflowExecutionJSON{
			Parameters: []Parameter{
				Parameter{
					Name: "vcopssuitepakBuildNumber",
					Type: "string",
					Value: Value{
						String: String{
							Value: build}}},
				Parameter{
					Name: "vcopssuitepakBuildType",
					Type: "string",
					Value: Value{
						String: String{
							Value: buildType}}},
				Parameter{
					Name: "email_to_address",
					Type: "string",
					Value: Value{
						String: String{
							Value: slackEmail}}},
				Parameter{
					Name: "vappPrefix",
					Type: "string",
					Value: Value{
						String: String{
							Value: slackUser}}},
				Parameter{
					Name: "vAppOwner",
					Type: "string",
					Value: Value{
						String: String{
							Value: slackUser}}}}}

		r, err := s.Post(href+"executions", &j, nil, nil)
		if err != nil {
			return err
		}
		if r.Status() != 202 {
			return fmt.Errorf("expect response status code 202. Actual %d", r.Status())
		}
	}

	text := comm.Mention(req.User)
	text += fmt.Sprintf("started vro workflow *%s* using *%s/%s*",
		wfName, build, buildType)

	iop := &comm.Interop{
		Request: req,
		Response: &comm.Response{
			Channel:    req.Channel,
			Text:       text,
			Parameters: comm.DefaultMessageParameters()}}
	comm.SB.ChanResponse <- iop.Response
	comm.SB.ChanPersist <- iop

	return nil
}
