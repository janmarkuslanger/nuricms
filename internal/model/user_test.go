package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserRoles(t *testing.T) {
	roles := GetUserRoles()
	assert.Len(t, roles, 2)
	assert.Contains(t, roles, RoleEditor)
	assert.Contains(t, roles, RoleAdmin)
}

func TestRoleConstants(t *testing.T) {
	var r Role
	r = RoleAdmin
	assert.Equal(t, "Admin", string(r))
	r = RoleEditor
	assert.Equal(t, "Editor", string(r))
}
