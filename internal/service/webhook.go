package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
)

type WebhookService interface {
	Create(name string, url string, requestType model.RequestType, events map[model.EventType]bool) (*model.Webhook, error)
	List(page, pageSize int) ([]model.Webhook, int64, error)
	FindByID(id uint) (*model.Webhook, error)
	Save(webhook *model.Webhook) error
	DeleteByID(id uint) error
	Dispatch(event string, payload any)
}

type webhookService struct {
	repos      *repository.Set
	httpClient *http.Client
}

func NewWebhookService(repos *repository.Set) WebhookService {
	return &webhookService{
		repos:      repos,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *webhookService) Create(name string, url string, requestType model.RequestType, events map[model.EventType]bool) (*model.Webhook, error) {
	var eventString strings.Builder
	for k, v := range events {
		if !v {
			continue
		}

		eventString.WriteString(string(k))
		eventString.WriteString(",")
	}

	webhook := model.Webhook{
		Name:        name,
		Url:         url,
		RequestType: requestType,
		Events:      eventString.String(),
	}

	err := s.repos.Webhook.Create(&webhook)
	return &webhook, err
}

func (s *webhookService) List(page, pageSize int) ([]model.Webhook, int64, error) {
	return s.repos.Webhook.List(page, pageSize)
}

func (s *webhookService) FindByID(id uint) (*model.Webhook, error) {
	return s.repos.Webhook.FindByID(id)
}

func (s *webhookService) Save(webhook *model.Webhook) error {
	return s.repos.Webhook.Save(webhook)
}

func (s *webhookService) DeleteByID(id uint) error {
	webhook, err := s.FindByID(id)
	if err != nil {
		return err
	}

	return s.repos.Webhook.Delete(webhook)
}

func (s *webhookService) Dispatch(event string, payload any) {
	hooks, err := s.repos.Webhook.ListByEvent(event)
	if err != nil {
		fmt.Println("Webhook find error:", err)
		return
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Webhook payload error:", err)
		return
	}

	for _, hook := range hooks {
		go func(h model.Webhook) {
			req, err := http.NewRequest(string(h.RequestType), h.Url, bytes.NewBuffer(body))
			if err != nil {
				fmt.Println("Webhook request error:", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := s.httpClient.Do(req)
			if err != nil {
				fmt.Println("Webhook delivery error:", err)
				return
			}
			resp.Body.Close()
		}(hook)
	}
}
