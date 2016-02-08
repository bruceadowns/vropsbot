package bot

import (
	"log"
	"reflect"
	"strings"

	"github.com/nlopes/slack"

	"github.com/bruceadowns/vropsbot/comm"
)

// Pump manages the message pump from slack
func Pump(rtm *slack.RTM, bot []string) (res int) {
	log.Print("Begin vRealize Operations Slack Bot Pump")

	// outer for loop label
LOOP:

	// infinite loop
	for {
		select {
		case event := <-rtm.IncomingEvents:
			// log.Printf("Incoming Event %s [%s]", reflect.TypeOf(event.Data), event.Data)
			switch message := event.Data.(type) {
			case *slack.MessageEvent:
				// handle messages @vropsbot

				c, u, t := message.Channel, message.User, message.Text
				if message.SubMessage != nil {
					u, t = message.SubMessage.User, message.SubMessage.Text
				}
				fields := strings.Fields(t)

				if fields[0] == bot[0] || fields[0] == bot[1] {
					err := PM.Dispatch(&comm.Request{
						Channel:   c,
						User:      u,
						Arguments: fields[1:]})
					if err != nil {
						log.Printf("Error dispatching: %s", err)
					}
				}
			case *slack.RTMError:
				log.Print("RTM Error", event)
			case *slack.InvalidAuthEvent:
				log.Print("Invalid authentication,", event)
				break LOOP
			case *slack.ConnectedEvent:
			case *slack.ConnectingEvent:
			case *slack.ConnectionErrorEvent:
			case *slack.DisconnectedEvent:
			case *slack.FilePublicEvent:
			case *slack.FileSharedEvent:
			case *slack.HelloEvent:
			case *slack.IMCreatedEvent:
			case *slack.LatencyReport:
			case *slack.PresenceChangeEvent:
			case *slack.TeamJoinEvent:
			case *slack.UserChangeEvent:
			case *slack.UserTypingEvent:
			default:
				log.Printf("Unexpected event %s [%s]", reflect.TypeOf(event), event)
			}
		}
	}

	log.Print("End vRealize Operations Slack Bot Pump")

	return
}
