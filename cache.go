package main

import (
	"errors"

	"github.com/nlopes/slack"
)

type cache struct {
	api      *slack.Client
	Users    map[string]slack.User
	Channels map[string]slack.Channel
	DMs      map[string]string
}

func NewCache(api *slack.Client) *cache {
	return &cache{
		api:      api,
		Users:    make(map[string]slack.User),
		Channels: make(map[string]slack.Channel),
		DMs:      make(map[string]string),
	}
}

func (c *cache) Populate() error {
	users, err := c.api.GetUsers()
	if err != nil {
		return errors.New("Couldn't pre-populate users")
	}

	userIDToName := make(map[string]string)
	for _, u := range users {
		c.Users[u.Name] = u
		userIDToName[u.ID] = u.Name
	}

	channels, err := c.api.GetChannels(true)
	if err != nil {
		return errors.New("Couldn't pre-populate channels")
	}
	for _, ch := range channels {
		c.Channels[ch.Name] = ch
	}

	dms, err := c.api.GetIMChannels()
	if err != nil {
		return errors.New("Couldn't pre-populate DMs")
	}
	for _, dm := range dms {
		name, ok := userIDToName[dm.User]
		if !ok {
			continue
		}
		c.DMs[name] = dm.ID
	}

	return nil
}
