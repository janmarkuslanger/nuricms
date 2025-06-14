package repository

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestWebhookRepository_ListByEvent(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewWebhookRepository(db)
	w1 := &model.Webhook{Url: "u1", Events: "a,b"}
	w2 := &model.Webhook{Url: "u2", Events: "b"}
	w3 := &model.Webhook{Url: "u3", Events: "c"}
	assert.NoError(t, repo.Create(w1))
	assert.NoError(t, repo.Create(w2))
	assert.NoError(t, repo.Create(w3))

	list, err := repo.ListByEvent("b")
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	ids := []uint{list[0].ID, list[1].ID}
	assert.Contains(t, ids, w1.ID)
	assert.Contains(t, ids, w2.ID)
}

func TestWebhookRepository_ListByEvent_None(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewWebhookRepository(db)
	w := &model.Webhook{Url: "u", Events: "x,y"}
	assert.NoError(t, repo.Create(w))

	list, err := repo.ListByEvent("z")
	assert.NoError(t, err)
	assert.Empty(t, list)
}
