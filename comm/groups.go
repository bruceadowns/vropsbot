package comm

import (
	"fmt"
	"log"

	"github.com/bruceadowns/vropsbot/common"

	"github.com/nlopes/slack"
)

// GroupCache provides a cache of slack groups
type GroupCache struct {
	client *slack.Client

	// map of slack ID:Name
	cache *common.BiMap
}

// single Groups instance
var Groups *GroupCache

// InitGroupCache is called by a main goroutine
// to prepopulate the group cache
func InitGroupCache(client *slack.Client) error {
	if Groups == nil {
		log.Printf("Initializing the group cache")

		Groups = &GroupCache{
			client: client,
			cache:  common.NewBiMap()}

		groups, err := client.GetGroups(true)
		if err != nil {
			return err
		}

		for _, v := range groups {
			Groups.cache.Put(v.ID, v.Name)
			SB.ChanPersist <- &common.SlackItem{
				ID:   v.ID,
				Name: v.Name,
				Type: common.SlackItemTypeGroup}
		}

		log.Printf("Done initializing the group cache")
	}

	return nil
}

// GetGroup returns a group given id
func GetGroup(id string) (string, error) {
	if Groups == nil {
		return "", fmt.Errorf("Group cache is not initialized")
	}

	if id == "" {
		return "", fmt.Errorf("Group id is empty")
	}

	name, exists := Groups.cache.GetByKey(id)
	if exists {
		return name, nil
	}

	g, err := Groups.client.GetGroupInfo(id)
	if err != nil {
		return "", err
	}
	if id != g.ID {
		return "", fmt.Errorf("given id '%s' does not match info id '%s'", id, g.ID)
	}

	log.Printf("Caching %s:%s\n", g.ID, g.Name)
	Groups.cache.Put(g.ID, g.Name)
	SB.ChanPersist <- &common.SlackItem{
		ID:   g.ID,
		Name: g.Name,
		Type: common.SlackItemTypeGroup}

	return g.Name, nil
}

// GetGroupByName returns a group given its name
func GetGroupByName(name string) (id string, err error) {
	if Groups == nil {
		err = fmt.Errorf("Group cache is not initialized")
		return
	}

	if name == "" {
		return "", fmt.Errorf("Group name is empty")
	}

	id, e := Groups.cache.GetByValue(name)
	if e {
		return id, nil
	}

	return "", fmt.Errorf("group %s not found", name)
}
