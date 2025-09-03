package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ApiClient struct {
	BaseURL    string
	ApiKey     string
	HTTPClient *http.Client
}

func New(baseURL, apiKey string) *ApiClient {
	return &ApiClient{
		BaseURL:    strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		ApiKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
	}
}

type ApiResponse[T any] struct {
	Success    bool        `json:"success"`
	Data       T           `json:"data"`
	Meta       *MetaData   `json:"meta,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type MetaData struct {
	Timestamp time.Time `json:"timestamp"`
}

type Pagination struct {
	PerPage int `json:"per_page"`
	Page    int `json:"page"`
}

type ContentItem struct {
	ID        uint                   `json:"id"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
	Values    map[string]interface{} `json:"values"`
}

func (c *ApiClient) get(path string, target any) error {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", c.ApiKey)
	req.Header.Set("Accept", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed: %s", res.Status)
	}

	return json.NewDecoder(res.Body).Decode(target)
}

func (c *ApiClient) FindContentByID(id uint) (*ContentItem, error) {
	var resp ApiResponse[*ContentItem]
	if err := c.get(fmt.Sprintf("/api/content/%d", id), &resp); err != nil {
		return nil, err
	}
	if !resp.Success || resp.Data == nil {
		return nil, errors.New("content not found")
	}
	return resp.Data, nil
}

func (c *ApiClient) FindContentByCollectionAlias(alias string, page, perPage int) ([]ContentItem, *Pagination, error) {
	path := fmt.Sprintf("/api/collections/%s/content?page=%d&perPage=%d", url.PathEscape(alias), page, perPage)
	var resp ApiResponse[[]ContentItem]
	if err := c.get(path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, resp.Pagination, nil
}

func (c *ApiClient) FindContentByCollectionAndFieldValue(alias, field, value string, page, perPage int) ([]ContentItem, *Pagination, error) {
	path := fmt.Sprintf("/api/collections/%s/content/filter?field=%s&value=%s&page=%d&perPage=%d",
		url.PathEscape(alias), url.QueryEscape(field), url.QueryEscape(value), page, perPage)
	var resp ApiResponse[[]ContentItem]
	if err := c.get(path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, resp.Pagination, nil
}
