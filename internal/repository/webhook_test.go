package repository

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupWebhookRepo(t *testing.T) *WebhookRepository {
	db, err := CreateTestDB()
	require.NoError(t, err)
	return NewWebhookRepository(db)
}

func TestWebhookRepository_CreateAndFindByID(t *testing.T) {
	repo := setupWebhookRepo(t)

	webhook := &model.Webhook{
		Url:    "https://example.com/hook",
		Events: "content:created,content:updated",
	}

	err := repo.Create(webhook)
	require.NoError(t, err)
	assert.NotZero(t, webhook.ID)

	found, err := repo.FindByID(webhook.ID)
	require.NoError(t, err)
	assert.Equal(t, webhook.Url, found.Url)
	assert.Equal(t, webhook.Events, found.Events)
}

func TestWebhookRepository_Save(t *testing.T) {
	repo := setupWebhookRepo(t)

	webhook := &model.Webhook{
		Url:    "https://example.com/save",
		Events: "content:deleted",
	}
	require.NoError(t, repo.Create(webhook))

	webhook.Url = "https://example.com/updated"
	err := repo.Save(webhook)
	require.NoError(t, err)

	updated, err := repo.FindByID(webhook.ID)
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/updated", updated.Url)
}

func TestWebhookRepository_Delete(t *testing.T) {
	repo := setupWebhookRepo(t)

	webhook := &model.Webhook{
		Url:    "https://example.com/delete",
		Events: "asset:created",
	}
	require.NoError(t, repo.Create(webhook))

	err := repo.Delete(webhook)
	require.NoError(t, err)

	_, err = repo.FindByID(webhook.ID)
	assert.Error(t, err)
}

func TestWebhookRepository_List(t *testing.T) {
	repo := setupWebhookRepo(t)

	for i := 0; i < 5; i++ {
		require.NoError(t, repo.Create(&model.Webhook{
			Url:    "https://example.com/" + string(rune('a'+i)),
			Events: "content:created",
		}))
	}

	list, total, err := repo.List(1, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, list, 2)

	list2, _, err := repo.List(3, 2)
	require.NoError(t, err)
	assert.Len(t, list2, 1)
}

func TestWebhookRepository_ListByEvent(t *testing.T) {
	repo := setupWebhookRepo(t)

	require.NoError(t, repo.Create(&model.Webhook{
		Url:    "https://example.com/a",
		Events: "content:created,content:updated",
	}))
	require.NoError(t, repo.Create(&model.Webhook{
		Url:    "https://example.com/b",
		Events: "asset:created",
	}))
	require.NoError(t, repo.Create(&model.Webhook{
		Url:    "https://example.com/c",
		Events: "content:created",
	}))

	hooks, err := repo.ListByEvent("content:created")
	require.NoError(t, err)
	assert.Len(t, hooks, 2)
}
