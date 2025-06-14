package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequestTypes(t *testing.T) {
	types := GetRequestTypes()
	assert.Len(t, types, 2)
	assert.Contains(t, types, RequestTypeGet)
	assert.Contains(t, types, RequestTypePost)
}

func TestGetWebhookEvents(t *testing.T) {
	events := GetWebhookEvents()
	assert.Len(t, events, 3)
	assert.Contains(t, events, EventContentCreated)
	assert.Contains(t, events, EventContentUpdated)
	assert.Contains(t, events, EventContentDeleted)
}

func TestRequestTypeConstants(t *testing.T) {
	var rt RequestType
	rt = RequestTypeGet
	assert.Equal(t, "GET", string(rt))
	rt = RequestTypePost
	assert.Equal(t, "POST", string(rt))
}

func TestEventTypeConstants(t *testing.T) {
	var et EventType
	et = EventContentCreated
	assert.Equal(t, "ContentCreated", string(et))
	et = EventContentUpdated
	assert.Equal(t, "ContentUpdated", string(et))
	et = EventContentDeleted
	assert.Equal(t, "ContentDeleted", string(et))
}
