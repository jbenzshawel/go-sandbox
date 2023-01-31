package role

import (
	"github.com/jbenzshawel/go-sandbox/common/cerror"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user/permission"
)

type Role struct {
	id   int
	name string

	permissions []*permission.Permission
}

func NewRole(name string, permissions []*permission.Permission) (*Role, error) {
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

func FromDatabase(roleID int, name string, permissions []*permission.Permission) (*Role, error) {
	r, err := NewRole(name, permissions)
	if err != nil {
		return nil, err
	}

	r.id = roleID

	return r, nil
}

func FromType(roleType Type) *Role {
	return &Role{
		id:   int(roleType),
		name: roleType.String(),
	}
}

func (r *Role) ID() int {
	return r.id
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) Type() Type {
	return Type(r.id)
}

func (r *Role) Permissions() []*permission.Permission {
	return r.permissions
}
