package repository

import (
	"strings"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository/base"
	"gorm.io/gorm"
)

type WebhookRepo interface {
	base.CRUDRepository[model.Webhook]
	ListByEvent(event string) ([]model.Webhook, error)
}

type webhookRepository struct {
	*base.BaseRepository[model.Webhook]
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) WebhookRepo {
	return &webhookRepository{
		BaseRepository: base.NewBaseRepository[model.Webhook](db),
		db:             db,
	}
}

func (r *webhookRepository) ListByEvent(event string) ([]model.Webhook, error) {
	var hooks []model.Webhook
	if err := r.db.Find(&hooks).Error; err != nil {
		return nil, err
	}
	var out []model.Webhook
	for _, h := range hooks {
		for _, ev := range strings.Split(h.Events, ",") {
			if ev == event {
				out = append(out, h)
				break
			}
		}
	}
	return out, nil
}
