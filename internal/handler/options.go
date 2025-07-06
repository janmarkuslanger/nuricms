package handler

type HandlerOptions struct {
	RedirectOnSuccess string
	RenderOnSuccess   string
	RenderOnFail      string
	RedirectOnFail    string
	TemplateData      func() (map[string]any, error)
}
