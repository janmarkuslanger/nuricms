package plugin

type HookPlugin interface {
	Name() string
	Register(h *HookRegistry)
}
