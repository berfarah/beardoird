package slack

import (
	"errors"

	"github.com/berfarah/beardroid/brain"
	"github.com/berfarah/beardroid/brain/redis"
	"github.com/nlopes/slack"
)

type User struct{ slack.User }

func (u User) ID() string        { return u.User.ID }
func (u User) Namespace() string { return "slack:user" }
func (u User) Indices() map[string]string {
	return map[string]string{"username": u.Name}
}

type Channel struct{ slack.Channel }

func (c Channel) ID() string        { return c.Channel.ID }
func (c Channel) Namespace() string { return "slack:channel" }
func (c Channel) Indices() map[string]string {
	return map[string]string{"name": c.Name}
}

type DM struct {
	Username     string
	UserID       string
	Conversation string
}

func (d DM) ID() string        { return d.Conversation }
func (d DM) Namespace() string { return "slack:dm" }
func (d DM) Indices() map[string]string {
	return map[string]string{"username": d.Username, "user_id": d.UserID}
}

type cache struct {
	api   *slack.Client
	brain *brain.Brain
}

func newCache(api *slack.Client) *cache {
	redis := redis.New("localhost:6379")
	return &cache{api: api, brain: brain.New(redis)}
}

func (c *cache) Fetch() error {
	users, err := c.api.GetUsers()
	if err != nil {
		return errors.New("Couldn't fetch users")
	}

	channels, err := c.api.GetChannels(true)
	if err != nil {
		return errors.New("Couldn't fetch channels")
	}

	dms, err := c.api.GetIMChannels()
	if err != nil {
		return errors.New("Couldn't fetch DMs")
	}

	return c.Populate(users, channels, dms)
}

func (c *cache) Populate(users []slack.User, channels []slack.Channel, dms []slack.IM) error {
	userIDToName := make(map[string]string)
	for _, u := range users {
		c.brain.Write(User{u})
		userIDToName[u.ID] = u.Name
	}

	for _, ch := range channels {
		c.brain.Write(Channel{ch})
	}

	for _, dm := range dms {
		if name, ok := userIDToName[dm.User]; ok {
			c.brain.Write(DM{
				Conversation: dm.ID,
				UserID:       dm.User,
				Username:     name,
			})
		}
	}

	return nil
}
