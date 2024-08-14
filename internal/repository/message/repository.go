package message

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	db "github.com/VadimGossip/platform_common/pkg/db/oracle"
	"github.com/sirupsen/logrus"

	"github.com/VadimGossip/mm_agent/internal/model"
	def "github.com/VadimGossip/mm_agent/internal/repository"
	"github.com/VadimGossip/mm_agent/internal/repository/message/converter"
	repoModel "github.com/VadimGossip/mm_agent/internal/repository/message/model"
)

const (
	messageTableName   string = "core_ext_msngrs_letters"
	msgIdColumn        string = "cextml_id"
	msgMsngrType       string = "cextml_msngr_type"
	msgSender          string = "cextml_sender"
	msgText            string = "cextml_text"
	msgStatus          string = "cextml_status"
	msgSendAt          string = "cextml_send_date"
	msgSendAfter       string = "cextml_send_after"
	addressesTableName string = "core_ext_mes_letter_addresses"
	addrAddressColumn  string = "cextmla_address"

	notSendStatus   int    = 0
	sendStatus      int    = 1
	mmMsgMsngrType  int    = 2
	teamNameTelejet string = "Telejet"
)

var _ def.MessageRepository = (*repository)(nil)

type repository struct {
	db db.Client
}

func NewRepository(db db.Client) *repository {
	return &repository{db: db}
}

func (r *repository) GetMessages(ctx context.Context) ([]*model.Message, error) {
	msgSelect := sq.Select(msgIdColumn, msgSender, msgText).
		From(messageTableName).
		PlaceholderFormat(sq.Colon).
		Where(sq.Eq{msgStatus: notSendStatus}).
		Where(sq.Eq{msgMsngrType: mmMsgMsngrType}).
		Where(sq.LtOrEq{msgSendAfter: time.Now()})

	query, args, err := msgSelect.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logrus.WithFields(logrus.Fields{
				"handler": "GetMessages",
				"problem": "rows close",
			}).Error(err)
		}
	}()

	var id int64
	var sender, text string
	msgs := make([]*model.Message, 0)
	for rows.Next() {
		if err = rows.Scan(&id, &sender, &text); err != nil {
			return nil, err
		}
		msgs = append(msgs, converter.ToMessageFromRepo(repoModel.Message{
			ID:     id,
			Sender: sender,
			Text:   text,
		}))
	}

	return msgs, nil
}

func (r *repository) GetMessageRecipients(ctx context.Context, id int64) ([]model.Recipient, error) {
	rcpSelect := sq.Select(addrAddressColumn).
		From(addressesTableName).
		PlaceholderFormat(sq.Colon).
		Where(sq.Eq{msgIdColumn: id})

	query, args, err := rcpSelect.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logrus.WithFields(logrus.Fields{
				"handler": "GetMessageRecipients",
				"problem": "rows close",
			}).Error(err)
		}
	}()

	var address string
	rcps := make([]model.Recipient, 0)
	for rows.Next() {
		if err = rows.Scan(&address); err != nil {
			return nil, err
		}
		rcps = append(rcps, converter.ToRecipientFromRepo(repoModel.Recipient{
			Name:     address,
			TeamName: teamNameTelejet,
			IsUser:   false,
		}))
	}

	return rcps, nil
}

func (r *repository) SetSendStatus(ctx context.Context, id int64, sendAt time.Time) error {
	setSendUpdate := sq.Update(messageTableName).
		PlaceholderFormat(sq.Colon).
		Set(msgSendAt, sendAt).
		Set(msgStatus, sendStatus).
		Where(sq.Eq{msgIdColumn: id})

	query, args, err := setSendUpdate.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
