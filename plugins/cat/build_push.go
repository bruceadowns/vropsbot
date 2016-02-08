package plugins

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bruceadowns/vropsbot/bot"
	"github.com/bruceadowns/vropsbot/comm"
	"github.com/bruceadowns/vropsbot/common"

	"gopkg.in/jmcvetta/napping.v3"
)

func init() {
	bot.Register(&buildPush{})
}

type buildPush struct {
	base
}

// Name returns build push's name
func (plugin *buildPush) Name() string {
	return "buildPush"
}

// Run starts the message pump for build push
func (plugin *buildPush) Run() (err error) {
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
		for _, t := range plugin.config.Targets {
			common.RandomDelay()

			go func(branch, target string) {
				pumpName := fmt.Sprintf("%s/%s/%s", plugin.Name(), branch, target)
				log.Printf("Starting pump for %s", pumpName)

				buildResult := unknownStatus
				for {
					select {
					case <-time.After(pollInterval):
						comm.SB.ChanHeartbeat <- pumpName

						pa := &common.PushActions{
							PluginName: plugin.Name()}

						url := plugin.config.URLs.BaseCatURL + plugin.config.URLs.BaseBuildURL

						params := DefaultCatParams()
						params["build__branch"] = branch
						params["targets"] = target

						p := params.AsUrlValues()
						s := napping.Session{}
						j := BuildJSON{}
						r, err := s.Get(url, &p, &j, nil)
						if err != nil {
							pa.Response = fmt.Sprintf("Error getting %s %s [%s]", url, params, err)
							pa.IsError = 1
							comm.SB.ChanPersist <- pa
							continue
						}
						if r.Status() != 200 {
							pa.Response = fmt.Sprintf("Expect response status code 200. Actual %d %s", r.Status(), url)
							pa.IsError = 1
							comm.SB.ChanPersist <- pa
							continue
						}
						if len(j.Objects) == 0 {
							// build not found
							continue
						}
						if len(j.Objects) > 1 {
							pa.Response = fmt.Sprintf("Expected one object. Actual %d %s", len(j.Objects), params)
							pa.IsError = 1
							comm.SB.ChanPersist <- pa
							continue
						}

						// if build status has changed
						//  and build status is not unknown
						//  and build is not the default status
						//  post the result
						if (buildResult != unknownStatus && buildResult != j.Objects[0].Result) ||
							(buildResult == unknownStatus && j.Objects[0].Result != defaultResult) {
							text := fmt.Sprintf("build *%s/%s* ", branch, target)
							text += fmt.Sprintf("changed *%s* to *%s* ", buildResult, j.Objects[0].Result)
							text += fmt.Sprintf("(%s)", comm.BuildwebSbURL(j.Objects[0].SbBuildID, strconv.Itoa(j.Objects[0].SbBuildID)))

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
			}(b, t)
		}
	}

	return
}
