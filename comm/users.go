package comm

import (
	"fmt"
	"log"

	"github.com/bruceadowns/vropsbot/common"

	"github.com/nlopes/slack"
)

// UserCache provides a cache of slack users
type UserCache struct {
	client *slack.Client

	// map of slack ID:Name
	cache *common.BiMap
}

// single Users instance
var Users *UserCache

// InitUserCache is called by a main goroutine
// to prepopulate the user cache
func InitUserCache(client *slack.Client) error {
	if Users == nil {
		log.Printf("Initializing the user cache")

		Users = &UserCache{
			client: client,
			cache:  common.NewBiMap()}

		users, err := client.GetUsers()
		if err != nil {
			return err
		}

		for _, v := range users {
			Users.cache.Put(v.ID, v.Name)
			SB.ChanPersist <- &common.SlackItem{
				ID:   v.ID,
				Name: v.Name,
				Type: common.SlackItemTypeUser}
		}

		log.Printf("Done initializing the user cache")
	}

	return nil
}

// GetUser returns a user given id
func GetUser(id string) (user string, err error) {
	if Users == nil {
		err = fmt.Errorf("User cache is not initialized")
		return
	}

	if id == "" {
		return "", fmt.Errorf("User id is empty")
	}

	user, ok := Users.cache.GetByKey(id)
	if ok {
		return
	}

	var u *slack.User
	u, err = Users.client.GetUserInfo(id)
	if err != nil {
		return
	}
	if id != u.ID {
		err = fmt.Errorf("given id '%s' does not match info id '%s'", id, u.ID)
		return
	}

	log.Printf("Caching %s:%s", u.ID, u.Name)
	Users.cache.Put(u.ID, u.Name)
	SB.ChanPersist <- &common.SlackItem{
		ID:   u.ID,
		Name: u.Name,
		Type: common.SlackItemTypeUser}

	user = u.ID

	return
}
