package app

import (
	"context"
	"fmt"

	"github.com/VadimGossip/platform_common/pkg/closer"
	"github.com/VadimGossip/platform_common/pkg/db/oracle"
	"github.com/VadimGossip/platform_common/pkg/db/oracle/odb"
	"github.com/VadimGossip/platform_common/pkg/db/oracle/transaction"
	"github.com/sirupsen/logrus"

	"github.com/VadimGossip/mm_agent/internal/client/mm"
	"github.com/VadimGossip/mm_agent/internal/model"
	"github.com/VadimGossip/mm_agent/internal/repository"
	msgRepo "github.com/VadimGossip/mm_agent/internal/repository/message"
	"github.com/VadimGossip/mm_agent/internal/service"
	msgService "github.com/VadimGossip/mm_agent/internal/service/message"
)

type serviceProvider struct {
	cfg *model.Config

	odbClient oracle.Client
	txManager oracle.TxManager
	msgRepo   repository.MessageRepository

	mmClient mm.Client

	msgService service.MessageService
}

func newServiceProvider(cfg *model.Config) *serviceProvider {
	return &serviceProvider{cfg: cfg}
}

func (s *serviceProvider) OdbClient(ctx context.Context) oracle.Client {
	if s.odbClient == nil {
		dsn := fmt.Sprintf(`user=%s password=%s connectString=%s:%d/%s`,
			s.cfg.OracleDb.Username,
			s.cfg.OracleDb.Password,
			s.cfg.OracleDb.Host,
			s.cfg.OracleDb.Port,
			s.cfg.OracleDb.Service)
		cl, err := odb.New(dsn)
		if err != nil {
			logrus.Fatalf("failed to create odb client: %s", err)
		}

		if err = cl.DB().Ping(ctx); err != nil {
			logrus.Fatalf("odb ping error: %s", err)
		}
		closer.Add(cl.Close)
		s.odbClient = cl
	}

	return s.odbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) oracle.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.OdbClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) MsgRepo(ctx context.Context) repository.MessageRepository {
	if s.msgRepo == nil {
		s.msgRepo = msgRepo.NewRepository(s.OdbClient(ctx))
	}
	return s.msgRepo
}

func (s *serviceProvider) MmClient() mm.Client {
	if s.mmClient == nil {
		mmClient, err := mm.NewClient(s.cfg.Mattermost.Url, s.cfg.Mattermost.Token)
		if err != nil {
			logrus.Fatalf("failed to authenticate mattermost bot: %s", err)
		}
		s.mmClient = mmClient
	}
	return s.mmClient
}

func (s *serviceProvider) MsgService(ctx context.Context) service.MessageService {
	if s.msgService == nil {
		s.msgService = msgService.NewService(s.MsgRepo(ctx), s.MmClient(), s.TxManager(ctx))
	}
	return s.msgService
}
