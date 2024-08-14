package mm

import (
	"fmt"
	"github.com/VadimGossip/mm_agent/internal/model"
	mm "github.com/mattermost/mattermost-server/v6/model"
	"time"
)

type Client interface {
	SendMessage(msg *model.Message) error
}

type client struct {
	cl    *mm.Client4
	botId string
}

func NewClient(url, pat string) (Client, error) {
	c := &client{cl: mm.NewAPIv4Client(url)}
	c.cl.SetOAuthToken(pat)
	user, _, err := c.cl.GetUser("me", "")
	if err != nil {
		return nil, err
	}
	c.botId = user.Id
	return c, nil
}

func (c *client) getChannelId(rcp model.Recipient) (string, error) {
	if rcp.IsUser {
		return "", fmt.Errorf("unsupported recipient")
	}

	team, _, err := c.cl.GetTeamByName(rcp.TeamName, "")
	if err != nil {
		return "", err
	}

	channel, _, err := c.cl.GetChannelByName(rcp.Name, team.Id, "")
	if err != nil {
		return "", err
	}
	return channel.Id, nil
}

func (c *client) sendMessageToChan(channelId, text string) error {
	_, _, err := c.cl.CreatePost(&mm.Post{
		ChannelId: channelId,
		Message:   text,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *client) SendMessage(msg *model.Message) error {
	for idx := range msg.Recipients {
		channelId, err := c.getChannelId(msg.Recipients[idx])
		if err != nil {
			return err
		}
		if err = c.sendMessageToChan(channelId, msg.Text); err != nil {
			return err
		}
		ts := time.Now()
		msg.SendAt = &ts
	}
	return nil
}
