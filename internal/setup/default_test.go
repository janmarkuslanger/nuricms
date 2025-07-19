package setup_test

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/setup"
	"github.com/janmarkuslanger/nuricms/pkg/config"
	"github.com/janmarkuslanger/nuricms/pkg/plugin"
	"github.com/stretchr/testify/assert"
)

func TestSetDefaultConfig_EmptyConfig(t *testing.T) {
	var opts config.Config
	conf := setup.SetDefaultConfig(opts)

	assert.Equal(t, conf.Port, "8080")
	var hooks []plugin.HookPlugin
	assert.Equal(t, conf.HookPlugins, hooks)
}
