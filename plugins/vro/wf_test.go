package plugins

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"gopkg.in/jmcvetta/napping.v3"
)

func TestVroUnauth(t *testing.T) {
	vro := "https://10.25.37.28:8281/vco/api/workflows/"
	user := "vropsbot"
	pass := "badpassword"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	s := napping.Session{
		Client:   &http.Client{Transport: tr},
		Userinfo: url.UserPassword(user, pass)}

	r, err := s.Get(vro, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.Status() != 401 {
		t.Fatal(fmt.Errorf("error calling vRO for [%d]", r.Status()))
	}

	t.Log("Successfully received 401 - unauthorized")
}

func TestVroUnsignedCert(t *testing.T) {
	vro := "https://10.25.37.28:8281/vco/api/workflows/"
	s := napping.Session{}
	_, err := s.Get(vro, nil, nil, nil)
	if err == nil {
		t.Fatal(err)
	}

	t.Logf("Successfully received bad certificate [%s]", err)
}

func TestWorkflows(t *testing.T) {
	vro := "https://10.25.37.28:8281/vco/api/workflows/"
	user := "vropsbot"
	pass := "vropsbot"

	wfname := "dev_checkin_tests"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	head := &http.Header{}
	head.Set("Accept-Type", "application/json")
	s := napping.Session{
		Client:   &http.Client{Transport: tr},
		Header:   head,
		Userinfo: url.UserPassword(user, pass)}

	np := napping.Params{
		"conditions": fmt.Sprintf("name=%s", wfname)}
	p := np.AsUrlValues()
	j := WorkflowsJSON{}
	r, err := s.Get(vro, &p, &j, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.Status() != 200 {
		t.Fatal(fmt.Errorf("error calling vRO for %s [%d]", wfname, r.Status()))
	}
	if len(j.Link) != 1 {
		t.Fatal(fmt.Errorf("expect one vRO %s workflow. Actual %d", wfname, len(j.Link)))
	}
	if len(j.Link[0].Href) == 0 {
		t.Fatal(fmt.Errorf("vRO workflow %s not found", wfname))
	}
	t.Logf("href: %s", j.Link[0].Href)
}

func TestWorkflowsUnknown(t *testing.T) {
	vro := "https://10.25.37.28:8281/vco/api/workflows/"
	user := "vropsbot"
	pass := "vropsbot"
	wfname := "unknown"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	head := &http.Header{}
	head.Set("Accept-Type", "application/json")
	s := napping.Session{
		Client:   &http.Client{Transport: tr},
		Header:   head,
		Userinfo: url.UserPassword(user, pass)}

	np := napping.Params{
		"conditions": fmt.Sprintf("name=%s", wfname)}
	p := np.AsUrlValues()
	j := WorkflowsJSON{}
	r, err := s.Get(vro, &p, &j, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.Status() != 200 {
		t.Fatal(fmt.Errorf("error calling vRO for %s [%d]", wfname, r.Status()))
	}
	if len(j.Link) != 0 {
		t.Fatal(fmt.Errorf("expect zero vRO %s workflow. Actual %d", wfname, len(j.Link)))
	}

	t.Log("Successfully not found 'unknown'")
}

func TestWorkflowPost(t *testing.T) {
	vro := "https://10.25.37.28:8281/vco/api/workflows/88e9de02-c974-4b0a-9c7d-30c75a8ec811/executions"
	user := "vropsbot"
	pass := "vropsbot"
	build := "6270899a"
	buildtype := "sba"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	head := &http.Header{}
	head.Set("Content-Type", "application/json")
	s := napping.Session{
		Client:   &http.Client{Transport: tr},
		Header:   head,
		Userinfo: url.UserPassword(user, pass),
		Log:      true}

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
						Value: buildtype}}},
			Parameter{
				Name: "email_to_address",
				Type: "string",
				Value: Value{
					String: String{
						Value: "vropsbot@vmware.com"}}},
			Parameter{
				Name: "vappPrefix",
				Type: "string",
				Value: Value{
					String: String{
						Value: "vropsbot"}}},
			Parameter{
				Name: "vAppOwner",
				Type: "string",
				Value: Value{
					String: String{
						Value: "vropsbot"}}}}}

	r, err := s.Post(vro, &j, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.Status() != 202 {
		t.Fatalf("Expect response status code 202. Actual %d", r.Status())
	}

	t.Log("Successfully posted an invalid sb")
}
