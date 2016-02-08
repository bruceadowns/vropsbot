package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nlopes/slack"

	"github.com/bruceadowns/vropsbot/bot"
	"github.com/bruceadowns/vropsbot/comm"
	"github.com/bruceadowns/vropsbot/common"

	_ "github.com/bruceadowns/vropsbot/plugins/admin"
	_ "github.com/bruceadowns/vropsbot/plugins/cat"
	_ "github.com/bruceadowns/vropsbot/plugins/vro"
)

func die(err error) {
	log.Fatal(err)
	os.Exit(-1)
}

func initCaches(client *slack.Client) {
	comm.InitUserCache(client)
	comm.InitChannelCache(client)
	comm.InitGroupCache(client)
}

func main() {
	config, err := common.NewConfig()
	if err != nil {
		die(fmt.Errorf("error reading configuration file: %v", err))
	}

	db, err := common.InitDatabase()
	if err != nil {
		die(fmt.Errorf("error initializing the database: %v", err))
	}
	defer db.Close()

	client := slack.New(config.SlackbotAuthToken)
	client.SetDebug(config.Debug)

	authRes, err := client.AuthTest()
	if err != nil {
		die(fmt.Errorf("error calling AuthTest: %s", err))
	}
	log.Printf("URL: %s", authRes.URL)
	log.Printf("Team: %s:%s", authRes.TeamID, authRes.Team)
	log.Printf("User: %s:%s", authRes.UserID, authRes.User)

	userID := authRes.UserID
	botIDs := []string{
		fmt.Sprintf("<@%s>:", userID),
		fmt.Sprintf("<@%s>", userID)}
	log.Printf("vropsbot identifiers: %s", botIDs)

	rtm := client.NewRTM()
	go rtm.ManageConnection()

	go initCaches(client)

	comm.SB.ChanResponse = comm.RespChan(client, config.SwitchBoard.ResponseChannelSize)
	comm.SB.ChanPersist = comm.SaveChan(config.SwitchBoard.SaveChannelSize, db)
	comm.SB.ChanHeartbeat = comm.HeartbeatChan(config.SwitchBoard.HeartbeatChannelSize, db)

	if err := bot.PM.Prepare(); err != nil {
		die(fmt.Errorf("error calling Plugin Manager Prepare: %s", err))
	}

	go bot.PM.Start()

	comm.SB.ChanPersist <- &common.KVStore{
		Key:   "start",
		Value: time.Now().String()}
	os.Exit(bot.Pump(rtm, botIDs))
}
