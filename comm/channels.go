package comm

import (
	"fmt"
	"log"

	"github.com/bruceadowns/vropsbot/common"
	"github.com/nlopes/slack"
)

// ChannelCache provides a cache of slack channels
type ChannelCache struct {
	client *slack.Client

	// map of slack ID:Name
	cache *common.BiMap
}

// single Channels instance
var Channels *ChannelCache

// InitChannelCache is called by a main goroutine
// to prepopulate the channel cache
func InitChannelCache(client *slack.Client) error {
	if Channels == nil {
		log.Printf("Initializing the channel cache")

		Channels = &ChannelCache{
			client: client,
			cache:  common.NewBiMap()}

		channels, err := client.GetChannels(true)
		if err != nil {
			return err
		}

		for _, v := range channels {
			Channels.cache.Put(v.ID, v.Name)
			SB.ChanPersist <- &common.SlackItem{
				ID:   v.ID,
				Name: v.Name,
				Type: common.SlackItemTypeChannel}
		}

		log.Printf("Done initializing the channel cache")
	}

	return nil
}

// GetChannel returns a channel given id
func GetChannel(id string) (string, error) {
	if Channels == nil {
		return "", fmt.Errorf("channel cache is not initialized")
	}

	if id == "" {
		return "", fmt.Errorf("channel id is empty")
	}

	channel, exists := Channels.cache.GetByKey(id)
	if exists {
		return channel, nil
	}

	c, err := Channels.client.GetChannelInfo(id)
	if err != nil {
		return "", err
	}
	if id != c.ID {
		return "", fmt.Errorf("given id '%s' does not match info id '%s'", id, c.ID)
	}

	log.Printf("Caching %s:%s\n", c.ID, c.Name)
	Channels.cache.Put(c.ID, c.Name)
	SB.ChanPersist <- &common.SlackItem{
		ID:   c.ID,
		Name: c.Name,
		Type: common.SlackItemTypeChannel}

	return c.Name, nil
}

// GetChannelByName returns a channel given its name
// Not sure how or whether this will be used
func GetChannelByName(name string) (string, error) {
	if Channels == nil {
		return "", fmt.Errorf("channel cache is not initialized")
	}

	if name == "" {
		return "", fmt.Errorf("channel name is empty")
	}

	id, exists := Channels.cache.GetByValue(name)
	if exists {
		return id, nil
	}

	return "", fmt.Errorf("channel %s not found", name)
}
