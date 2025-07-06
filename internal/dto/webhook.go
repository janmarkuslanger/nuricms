package dto

type WebhookData struct {
	Name        string
	Url         string
	RequestType string
	Events      map[string]bool
}
