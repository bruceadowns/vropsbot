package plugins

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"gopkg.in/jmcvetta/napping.v3"
)

func catURL(t *testing.T) (res string) {
	res = os.Getenv("CAT_URL")
	if res == "" {
		t.Skip("Provide CAT_URL. i.e. CAT_URL=http://cat.eng.vmware.com")
	}

	return
}

func TestUnmarshal(t *testing.T) {
	var j bytes.Buffer
	j.WriteString("{}")
	var b BuildJSON

	d := json.Unmarshal(j.Bytes(), &b)
	t.Log(d)
}

func TestCatBuildVcopssuitevmMainPass(t *testing.T) {
	url := catURL(t)
	url += "/api/v2.0/deliverable"
	t.Log("CAT URL:", url)

	params := DefaultCatParams()
	params["targets"] = "vcopssuitevm"
	params["build__branch"] = "main"
	params["result"] = "PASS"
	p := params.AsUrlValues()
	t.Logf("p: %v\n", p)

	s := napping.Session{}
	j := BuildJSON{}
	response, err := s.Get(url, &p, &j, nil)
	if err != nil {
		t.Fatal(err)
	}
	if response.Status() != 200 {
		t.Fatalf("Bad response status code: %d\n", response.Status())
	}

	if len(j.Objects) != 1 {
		t.Fatalf("Expected one object. Actual %d\n", len(j.Objects))
	}
	t.Logf("sbbuildid: %d\n", j.Objects[0].SbBuildID)
	t.Logf("result: %s\n", j.Objects[0].Result)
}
