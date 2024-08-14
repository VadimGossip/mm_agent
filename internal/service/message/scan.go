package message

import (
	"context"
	"github.com/VadimGossip/mm_agent/internal/model"
)

func (s *service) scan(ctx context.Context) ([]*model.Message, error) {
	msgs := make([]*model.Message, 0)
	err := s.txManager.ReadSerializable(ctx, func(ctx context.Context) error {
		var txErr error
		msgs, txErr = s.msgRepository.GetMessages(ctx)
		if txErr != nil {
			return txErr
		}

		for idx := range msgs {
			msgs[idx].Recipients, txErr = s.msgRepository.GetMessageRecipients(ctx, msgs[idx].ID)
			if txErr != nil {
				return txErr
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
