package repository

import (
	"errors"
	"strings"

	"github.com/janmarkuslanger/nuricms/model"
	"gorm.io/gorm"
)

type WebhookRepository struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) *WebhookRepository {
	return &WebhookRepository{
		db: db,
	}
}

func (r *WebhookRepository) Save(webhook *model.Webhook) error {
	if err := r.db.Save(webhook).Error; err != nil {
		return err
	}

	return nil
}

func (r *WebhookRepository) Create(webhook *model.Webhook) error {
	return r.db.Create(webhook).Error
}

func (r *WebhookRepository) Delete(webhook *model.Webhook) error {
	err := r.db.Delete(webhook).Error
	return err
}

func (r *WebhookRepository) FindByID(id uint) (*model.Webhook, error) {
	var webhook model.Webhook
	err := r.db.Where("id = ?", id).First(&webhook).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &webhook, err
}

func (r *WebhookRepository) List(page, pageSize int) ([]model.Webhook, int64, error) {
	var webhooks []model.Webhook
	var totalCount int64

	err := r.db.Model(&model.Webhook{}).
		Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err = r.db.
		Offset(offset).
		Limit(pageSize).
		Find(&webhooks).Error

	return webhooks, totalCount, err
}

func (r *WebhookRepository) ListByEvent(event string) ([]model.Webhook, error) {
	var webhooks []model.Webhook
	if err := r.db.Find(&webhooks).Error; err != nil {
		return nil, err
	}

	var webhooksByEvent []model.Webhook

	for _, hook := range webhooks {
		for _, ev := range strings.Split(hook.Events, ",") {
			if ev == event {
				webhooksByEvent = append(webhooksByEvent, hook)
				break
			}
		}
	}

	return webhooksByEvent, nil
}
