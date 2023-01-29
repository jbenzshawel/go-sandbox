package user

import (
	"github.com/jbenzshawel/go-sandbox/common/cerror"
)

type Role struct {
	id   int
	name string

	permissions []*Permission
}

func NewRole(name string, permissions []*Permission) (*Role, error) {
	validationErrors := map[string]string{}

	if name == "" {
		validationErrors["name"] = "name cannot be empty"
	}

	if len(validationErrors) > 0 {
		return nil, cerror.NewValidationError("invalid role", validationErrors)
	}

	return &Role{
		name:        name,
		permissions: permissions,
	}, nil
}

func RoleFromDatabase(roleID int, name string, permissions []*Permission) (*Role, error) {
	r, err := NewRole(name, permissions)
	if err != nil {
		return nil, err
	}

	r.id = roleID

	return r, nil
}

func (r *Role) ID() int {
	return r.id
}

func (r *Role) Name() string {
	return r.name
}
