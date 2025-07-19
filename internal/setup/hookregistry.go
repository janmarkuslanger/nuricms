package setup

import "github.com/janmarkuslanger/nuricms/pkg/plugin"

func InitHookRegistry(pgs []plugin.HookPlugin) (hr *plugin.HookRegistry) {
	hr = plugin.NewHookRegistry()
	for _, p := range pgs {
		if p == nil {
			continue
		}
		p.Register(hr)
	}

	return hr
}
