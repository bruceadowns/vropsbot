package plugins

import (
	"fmt"
	"log"
	"time"

	"github.com/bruceadowns/vropsbot/bot"
	"github.com/bruceadowns/vropsbot/comm"
	"github.com/bruceadowns/vropsbot/common"

	"gopkg.in/jmcvetta/napping.v3"
)

func init() {
	bot.Register(&testrunPush{})
}

type testrunPush struct {
	base
}

// Name returns testrunPush name
func (plugin *testrunPush) Name() string {
	return "testrunPush"
}

// Run testrunPush
func (plugin *testrunPush) Run() (err error) {
	pollInterval := common.DefaultPollInterval
	var channels []string
	var groups []string

	pollInterval, err = time.ParseDuration(plugin.config.PushPluginConfig.PollInterval)
	if err != nil {
		return err
	}
	channels = plugin.config.PushPluginConfig.Channels
	groups = plugin.config.PushPluginConfig.Groups

	if len(channels) == 0 && len(groups) == 0 {
		return fmt.Errorf("no channels or groups specified")
	}

	log.Printf("%s poll interval %d", plugin.Name(), pollInterval)
	log.Printf("%s channels %v", plugin.Name(), channels)
	log.Printf("%s groups %v", plugin.Name(), groups)

	for _, b := range plugin.config.Branches {
		for _, a := range plugin.config.Areas {
			common.RandomDelay()

			go func(branch, area string) {
				pumpName := fmt.Sprintf("%s/%s/%s", plugin.Name(), branch, area)
				log.Printf("Starting pump for %s", pumpName)

				buildResult := unknownStatus
				for {
					select {
					case <-time.After(pollInterval):
						comm.SB.ChanHeartbeat <- pumpName

						pa := &common.PushActions{
							PluginName: plugin.Name()}

						url := plugin.config.URLs.BaseCatURL + plugin.config.URLs.BaseTestrunURL

						params := DefaultCatParams()
						params["deliverables__build__branch"] = branch
						params["area__name"] = area

						p := params.AsUrlValues()
						s := napping.Session{}
						j := TestrunJSON{}
						r, err := s.Get(url, &p, &j, nil)
						if err != nil {
							pa.Response = fmt.Sprintf("Error getting %s %s [%s]", url, params, err)
							pa.IsError = 1
							comm.SB.ChanPersist <- pa
							continue
						}
						if r.Status() != 200 {
							pa.Response = fmt.Sprintf("Expect response status code 200. Actual %d. %s", r.Status(), params)
							pa.IsError = 1
							comm.SB.ChanPersist <- pa
							continue
						}
						if len(j.Objects) == 0 {
							// build not found
							continue
						}
						if len(j.Objects) > 1 {
							pa.Response = fmt.Sprintf("Expected one object. Actual %d %s\n", len(j.Objects), params)
							pa.IsError = 1
							comm.SB.ChanPersist <- pa
							continue
						}

						// if first time and build result is not PASS
						//  or build result has changed
						//  post results
						if (buildResult != unknownStatus && buildResult != j.Objects[0].Result) ||
							(buildResult == unknownStatus && j.Objects[0].Result != defaultResult) {
							text := fmt.Sprintf("testrun %s ",
								comm.CatTestrunURL(j.Objects[0].ID, fmt.Sprintf("%s/%s", branch, area)))
							text += fmt.Sprintf("changed *%s* to *%s* ", buildResult, j.Objects[0].Result)
							text += fmt.Sprintf("(%s)", comm.CatTestrunResultsURL(j.Objects[0].ResultsDir))

							pa.Response = text
							pa.Before = buildResult
							pa.After = j.Objects[0].Result

							for _, c := range channels {
								cname, err := comm.GetChannelByName(c)
								if err != nil {
									log.Print(err)
									continue
								}
								comm.SB.ChanResponse <- &comm.Response{
									Channel:    cname,
									Text:       text,
									Parameters: comm.DefaultMessageParameters()}

								pa.Channel = cname
								comm.SB.ChanPersist <- pa
							}
							for _, g := range groups {
								gname, err := comm.GetGroupByName(g)
								if err != nil {
									log.Print(err)
									continue
								}

								comm.SB.ChanResponse <- &comm.Response{
									Channel:    gname,
									Text:       text,
									Parameters: comm.DefaultMessageParameters()}

								pa.Channel = gname
								comm.SB.ChanPersist <- pa
							}
						}

						buildResult = j.Objects[0].Result
					}
				}
			}(b, a)
		}
	}

	return
}
