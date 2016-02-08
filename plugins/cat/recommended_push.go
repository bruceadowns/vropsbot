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
	bot.Register(&recommendedPush{})
}

type recommendedPush struct {
	base
}

// Name returns recommendedPush name
func (plugin *recommendedPush) Name() string {
	return "recommendedPush"
}

// Run recommendedPush
func (plugin *recommendedPush) Run() (err error) {
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
		common.RandomDelay()

		go func(branch string) {
			pumpName := fmt.Sprintf("%s/%s", plugin.Name(), branch)
			log.Printf("Starting pump for %s", pumpName)

			currentChangeset := unknownStatus
			for {
				select {
				case <-time.After(pollInterval):
					comm.SB.ChanHeartbeat <- pumpName

					pa := &common.PushActions{
						PluginName: plugin.Name()}

					jRecommended := RecommendedJSON{}
					{
						url := plugin.config.URLs.BaseCatURL + plugin.config.URLs.BaseRecommendedURL
						params := DefaultRecommendedParams()
						params["branch__name"] = branch

						p := params.AsUrlValues()
						s := napping.Session{}
						r, err := s.Get(url, &p, &jRecommended, nil)
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
						if len(jRecommended.Objects) == 0 {
							// build not found
							continue
						}
						if len(jRecommended.Objects) > 1 {
							pa.Response = fmt.Sprintf("Expected one object. Actual %d %s\n", len(jRecommended.Objects), params)
							pa.IsError = 1
							comm.SB.ChanPersist <- pa
							continue
						}
					}

					jCurrBuild := CurrBuildJSON{}
					{
						url := plugin.config.URLs.BaseCatURL + jRecommended.Objects[0].CurrBuild
						params := DefaultCurrBuildParams()

						p := params.AsUrlValues()
						s := napping.Session{}
						r, err := s.Get(url, &p, &jCurrBuild, nil)
						if err != nil {
							pa.Response = fmt.Sprintf("Error getting %s [%s]", url, err)
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
					}

					if currentChangeset != unknownStatus &&
						currentChangeset != jCurrBuild.Changeset {
						text := fmt.Sprintf("recommendation for *%s* ", branch)
						text += fmt.Sprintf(
							"changed from %s to %s",
							comm.P4WebBranchURL(branch, currentChangeset),
							comm.P4WebBranchURL(branch, jCurrBuild.Changeset))

						pa.Response = text
						pa.Before = currentChangeset
						pa.After = jCurrBuild.Changeset

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

					currentChangeset = jCurrBuild.Changeset
				}
			}
		}(b)
	}

	return
}
