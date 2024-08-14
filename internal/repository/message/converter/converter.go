package converter

import (
	"github.com/VadimGossip/mm_agent/internal/model"
	repoModel "github.com/VadimGossip/mm_agent/internal/repository/message/model"
)

func ToMessageFromRepo(msg repoModel.Message) *model.Message {
	return &model.Message{
		ID:     msg.ID,
		Sender: msg.Sender,
		Text:   msg.Text,
	}
}

func ToRecipientFromRepo(rcp repoModel.Recipient) model.Recipient {
	return model.Recipient{
		Name:     rcp.Name,
		TeamName: rcp.TeamName,
		IsUser:   rcp.IsUser,
	}
}
