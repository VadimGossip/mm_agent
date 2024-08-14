package repository

import (
	"context"
	"github.com/VadimGossip/mm_agent/internal/model"
	"time"
)

type MessageRepository interface {
	GetMessages(ctx context.Context) ([]*model.Message, error)
	GetMessageRecipients(ctx context.Context, id int64) ([]model.Recipient, error)
	SetSendStatus(ctx context.Context, id int64, sendAt time.Time) error
}
