package message

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *service) RunScanner(ctx context.Context) {
	logrus.Info("Mattermost message scanner started")
	defer logrus.Info("Mattermost message scanner started")
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			msgs, err := s.scan(ctx)
			if err != nil {
				logrus.Errorf("Mattermost message scanner scan err: %s", err)
				continue
			}

			if len(msgs) > 0 {
				if err = s.send(ctx, msgs); err != nil {
					logrus.Errorf("Mattermost message scanner send err: %s", err)
				}
			}
		}
	}
}
