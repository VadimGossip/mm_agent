package message

import (
	"context"
	"github.com/VadimGossip/mm_agent/internal/model"
)

func (s *service) send(ctx context.Context, msgs []*model.Message) error {
	return s.txManager.ReadSerializable(ctx, func(ctx context.Context) error {
		for idx := range msgs {
			if txErr := s.mmClient.SendMessage(msgs[idx]); txErr != nil {
				return txErr
			}
			if txErr := s.msgRepository.SetSendStatus(ctx, msgs[idx].ID, *msgs[idx].SendAt); txErr != nil {
				return txErr
			}
		}
		return nil
	})
}
