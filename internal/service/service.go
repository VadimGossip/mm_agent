package service

import "context"

type MessageService interface {
	RunScanner(ctx context.Context)
}
