package model

import (
	"gorm.io/gorm"
)

type EventType string

const (
	EventContentCreated EventType = "ContentCreated"
	EventContentUpdated EventType = "ContentUpdated"
	EventContentDeleted EventType = "ContentDeleted"
)

type RequestType string

const (
	RequestTypeGet  RequestType = "GET"
	RequestTypePost RequestType = "POST"
)

type Webhook struct {
	gorm.Model
	Name        string      `gorm:"size:80;not null"`
	Url         string      `gorm:"not null"`
	RequestType RequestType `gorm:"size:80;not null"`
	Active      bool        `gorm:"not null;default:true"`
	Events      string      `gorm:"type:varchar(500);not null"`
}

func GetRequestTypes() []RequestType {
	return []RequestType{
		RequestTypeGet,
		RequestTypePost,
	}
}

func GetWebhookEvents() []EventType {
	return []EventType{
		EventContentCreated,
		EventContentDeleted,
		EventContentUpdated,
	}
}
