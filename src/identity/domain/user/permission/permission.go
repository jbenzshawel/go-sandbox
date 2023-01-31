package permission

import (
	"github.com/jbenzshawel/go-sandbox/common/cerror"
)

type Permission struct {
	id   int
	name string
}

func NewPermission(name string) (*Permission, error) {
	validationErrors := map[string]string{}
	if name == "" {
		validationErrors["name"] = "name cannot be empty"
	}

	if len(validationErrors) > 0 {
		return nil, cerror.NewValidationError("invalid role", validationErrors)
	}

	return &Permission{
		name: name,
	}, nil
}

func FromDatabase(permissionID int, name string) (*Permission, error) {
	p, err := NewPermission(name)
	if err != nil {
		return nil, err
	}

	p.id = permissionID

	return p, nil
}

func (p *Permission) ID() int {
	return p.id
}

func (p *Permission) Name() string {
	return p.name
}

func (p *Permission) Type() Type {
	return Type(p.id)
}
