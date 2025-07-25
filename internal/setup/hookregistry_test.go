package setup_test

import (
	"errors"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/setup"
	"github.com/janmarkuslanger/nuricms/pkg/plugin"
	"github.com/stretchr/testify/assert"
)

type mockPlugin struct {
	called *bool
}

func (m *mockPlugin) Name() string {
	return "mock"
}

func (m *mockPlugin) Register(h *plugin.HookRegistry) {
	*m.called = true
	h.Register("test", func(payload any) error {
		return nil
	})
}

func TestInitHookRegistry_AllPluginsRegistered(t *testing.T) {
	wasCalled := false
	p := &mockPlugin{called: &wasCalled}
	registry := setup.InitHookRegistry([]plugin.HookPlugin{p})

	assert.NotNil(t, registry)
	assert.True(t, wasCalled, "plugin.Register should have been called")

	err := registry.Run("test", nil)
	assert.NoError(t, err, "hook should have executed successfully")
}

func TestInitHookRegistry_IgnoresNilPlugin(t *testing.T) {
	registry := setup.InitHookRegistry([]plugin.HookPlugin{nil})

	assert.NotNil(t, registry)
	err := registry.Run("any", nil)
	assert.NoError(t, err)
}

func TestInitHookRegistry_MultipleHooks(t *testing.T) {
	registry := plugin.NewHookRegistry()

	count := 0
	registry.Register("multi", func(payload any) error {
		count++
		return nil
	})
	registry.Register("multi", func(payload any) error {
		count++
		return nil
	})

	err := registry.Run("multi", nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestInitHookRegistry_HookReturnsError(t *testing.T) {
	registry := plugin.NewHookRegistry()

	registry.Register("fail", func(payload any) error {
		return errors.New("hook failed")
	})

	err := registry.Run("fail", nil)
	assert.Error(t, err)
	assert.EqualError(t, err, "hook failed")
}
