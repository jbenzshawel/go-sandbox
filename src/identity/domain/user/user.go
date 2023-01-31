package user

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user/permission"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user/role"
)

type User struct {
	id   int
	uuid uuid.UUID

	firstName     string
	lastName      string
	email         string
	emailVerified bool
	enabled       bool

	roles       []*role.Role
	permissions map[permission.Type]struct{}

	createdAt     time.Time
	lastUpdatedAt time.Time
}

func NewUser(firstName, lastName, email string, emailVerified, enabled bool) (*User, error) {
	validationErrors := map[string]string{}
	if firstName == "" {
		validationErrors["firstName"] = "firstName cannot be empty"
	}
	if lastName == "" {
		validationErrors["lastName"] = "lastName cannot be empty"
	}
	if email == "" {
		validationErrors["email"] = "email cannot be empty"
	}

	if len(validationErrors) > 0 {
		return nil, cerror.NewValidationError("invalid user", validationErrors)
	}

	return &User{
		firstName:     firstName,
		lastName:      lastName,
		email:         email,
		emailVerified: emailVerified,
		enabled:       enabled,
		createdAt:     time.Now().UTC(),
		lastUpdatedAt: time.Now().UTC(),
	}, nil
}

func FromDatabase(id int, userUUID uuid.UUID, firstName, lastName, email string,
	emailVerified, enabled bool, roles []*role.Role, createdAt, lastUpdatedAt time.Time) (*User, error) {
	u, err := NewUser(firstName, lastName, email, emailVerified, enabled)
	if err != nil {
		return nil, err
	}

	u.id = id
	u.uuid = userUUID
	u.roles = roles
	u.setPermissions(roles)
	u.createdAt = createdAt
	u.lastUpdatedAt = lastUpdatedAt

	return u, nil
}

func (u *User) ID() int {
	return u.id
}

func (u *User) UUID() uuid.UUID {
	return u.uuid
}

func (u *User) FirstName() string {
	return u.firstName
}

func (u *User) LastName() string {
	return u.lastName
}

func (u *User) Email() string {
	return u.email
}

func (u *User) EmailVerified() bool {
	return u.emailVerified
}

func (u *User) Enabled() bool {
	return u.enabled
}

func (u *User) Roles() []*role.Role {
	return u.roles
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) LastUpdatedAt() time.Time {
	return u.lastUpdatedAt
}

func (u *User) SetID(id int) error {
	if id < 1 {
		return errors.New("user id must be greater than 0")
	}
	u.id = id
	u.lastUpdatedAt = time.Now().UTC()
	return nil
}

func (u *User) SetUUID(userUUID uuid.UUID) error {
	if userUUID == uuid.Nil {
		return errors.New("user uuid cannot be nil")
	}
	u.uuid = userUUID
	u.lastUpdatedAt = time.Now().UTC()
	return nil
}

func (u *User) SetEmailVerified(verified bool) {
	u.emailVerified = verified
	u.lastUpdatedAt = time.Now().UTC()
}

func (u *User) AddRole(roleType role.Type) {
	u.roles = append(u.roles, role.FromType(roleType))
	u.lastUpdatedAt = time.Now().UTC()
}

func (u *User) setRoles(roles []*role.Role) {
	u.roles = roles
	u.setPermissions(roles)
}

func (u *User) setPermissions(roles []*role.Role) {
	if u.permissions == nil {
		u.permissions = map[permission.Type]struct{}{}
	}

	for _, r := range roles {
		for _, p := range r.Permissions() {
			u.permissions[p.Type()] = struct{}{}
		}
	}
}

func (u *User) HasPermission(permitType permission.Type) bool {
	_, ok := u.permissions[permitType]
	return ok
}
