package message

import (
	"github.com/VadimGossip/mm_agent/internal/client/mm"
	"github.com/VadimGossip/mm_agent/internal/repository"
	db "github.com/VadimGossip/platform_common/pkg/db/oracle"

	def "github.com/VadimGossip/mm_agent/internal/service"
)

var _ def.MessageService = (*service)(nil)

type service struct {
	msgRepository repository.MessageRepository
	mmClient      mm.Client
	txManager     db.TxManager
}

func NewService(msgRepository repository.MessageRepository,
	mmClient mm.Client,
	txManager db.TxManager) *service {
	return &service{
		msgRepository: msgRepository,
		mmClient:      mmClient,
		txManager:     txManager,
	}
}
