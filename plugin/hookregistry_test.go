package plugin

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHookRegistry_RegisterAndRun(t *testing.T) {
	reg := NewHookRegistry()

	var called bool
	testPayload := "data"

	reg.Register("test:event", func(p any) error {
		assert.Equal(t, testPayload, p)
		called = true
		return nil
	})

	err := reg.Run("test:event", testPayload)

	assert.NoError(t, err)
	assert.True(t, called, "Hook function should have been called")
}

func TestHookRegistry_MultipleHooks(t *testing.T) {
	reg := NewHookRegistry()

	callOrder := []int{}

	reg.Register("test:event", func(p any) error {
		callOrder = append(callOrder, 1)
		return nil
	})
	reg.Register("test:event", func(p any) error {
		callOrder = append(callOrder, 2)
		return nil
	})

	err := reg.Run("test:event", nil)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2}, callOrder)
}

func TestHookRegistry_HookReturnsError(t *testing.T) {
	reg := NewHookRegistry()

	reg.Register("test:event", func(p any) error {
		return errors.New("fail")
	})
	reg.Register("test:event", func(p any) error {
		t.Fatal("This should not be called after failure")
		return nil
	})

	err := reg.Run("test:event", nil)

	assert.Error(t, err)
	assert.EqualError(t, err, "fail")
}
