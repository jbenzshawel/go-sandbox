package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
)

type Role struct {
	id   int
	uuid uuid.UUID

	name        string
	description string

	permissions []*Permission

	createdAt     time.Time
	lastUpdatedAt time.Time
}

func NewRole(roleUUID uuid.UUID, name, description string, permissions []*Permission) (*Role, error) {
	validationErrors := map[string]string{}
	if roleUUID == uuid.Nil {
		validationErrors["uuid"] = "UUID cannot be empty"
	}
	if name == "" {
		validationErrors["name"] = "name cannot be empty"
	}

	if len(validationErrors) > 0 {
		return nil, cerror.NewValidationError("invalid role", validationErrors)
	}

	return &Role{
		uuid:          roleUUID,
		name:          name,
		description:   description,
		permissions:   permissions,
		createdAt:     time.Now().UTC(),
		lastUpdatedAt: time.Now().UTC(),
	}, nil
}

func RoleFromDatabase(roleID int, roleUUID uuid.UUID, name, description string,
	permissions []*Permission, createdAt, lastUpdatedAt time.Time) (*Role, error) {
	r, err := NewRole(roleUUID, name, description, permissions)
	if err != nil {
		return nil, err
	}

	r.id = roleID
	r.createdAt = createdAt
	r.lastUpdatedAt = lastUpdatedAt

	return r, nil
}

func (r *Role) ID() int {
	return r.id
}

func (r *Role) UUID() uuid.UUID {
	return r.uuid
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) Description() string {
	return r.description
}

func (r *Role) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Role) LastUpdatedAt() time.Time {
	return r.lastUpdatedAt
}
