package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jbenzshawel/go-sandbox/identity/domain/user/permission"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user/role"
)

func TestHasPermission(t *testing.T) {
	u, err := NewUser("John", "Smith", "john.smith@email.com", true, true)
	require.NoError(t, err)
	assert.False(t, u.HasPermission(permission.EditUsers))

	p, err := permission.FromDatabase(int(permission.EditUsers), "Edit Users")
	require.NoError(t, err)
	r, err := role.FromDatabase(int(role.Admin), "Admin", []*permission.Permission{p})
	u.setRoles([]*role.Role{r})
	assert.True(t, u.HasPermission(permission.EditUsers))
}
