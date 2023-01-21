package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
)

type Permission struct {
	id   int
	uuid uuid.UUID

	name        string
	description string

	createdAt     time.Time
	lastUpdatedAt time.Time
}

func NewPermission(roleUUID uuid.UUID, name, description string) (*Permission, error) {
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

	return &Permission{
		uuid:          roleUUID,
		name:          name,
		description:   description,
		createdAt:     time.Now().UTC(),
		lastUpdatedAt: time.Now().UTC(),
	}, nil
}

func PermissionFromDatabase(permissionID int, permissionUUID uuid.UUID, name, description string, createdAt, lastUpdatedAt time.Time) (*Permission, error) {
	p, err := NewPermission(permissionUUID, name, description)
	if err != nil {
		return nil, err
	}

	p.id = permissionID
	p.createdAt = createdAt
	p.lastUpdatedAt = lastUpdatedAt

	return p, nil
}

func (p *Permission) ID() int {
	return p.id
}

func (p *Permission) UUID() uuid.UUID {
	return p.uuid
}

func (p *Permission) Name() string {
	return p.name
}

func (p *Permission) Description() string {
	return p.description
}

func (p *Permission) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Permission) LastUpdatedAt() time.Time {
	return p.lastUpdatedAt
}
