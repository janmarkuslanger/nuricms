package plugin

type HookFunc func(payload any) error

type HookRegistry struct {
	hooks map[string][]HookFunc
}

func NewHookRegistry() *HookRegistry {
	return &HookRegistry{
		hooks: make(map[string][]HookFunc),
	}
}

func (hr *HookRegistry) Register(name string, fn HookFunc) {
	hr.hooks[name] = append(hr.hooks[name], fn)
}

func (hr *HookRegistry) Run(name string, payload any) error {
	for _, fn := range hr.hooks[name] {
		if err := fn(payload); err != nil {
			return err
		}
	}
	return nil
}
