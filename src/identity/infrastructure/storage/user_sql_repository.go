package storage

import (
	"database/sql"
	"os"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/common/database"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/gen/identity/identity/model"
	. "github.com/jbenzshawel/go-sandbox/identity/infrastructure/gen/identity/identity/table"
)

type UserSqlRepository struct {
	dbProvider database.DbProvider
}

func NewUserSqlRepository(dbProvider database.DbProvider) *UserSqlRepository {
	return &UserSqlRepository{
		dbProvider: dbProvider,
	}
}

func TryCreateUserSqlRepository() (*UserSqlRepository, bool) {
	if connectionString, ok := os.LookupEnv("IDENTITY_POSTGRES"); ok {
		return NewUserSqlRepository(func() (*sql.DB, error) {
			return sql.Open("postgres", connectionString)
		}), true
	}

	return nil, false
}

func (r *UserSqlRepository) CreateUser(u *user.User) error {
	_, err := database.Execute(r.dbProvider, Users.INSERT(Users.MutableColumns).
		MODEL(model.Users{
			UserUUID:      u.UUID(),
			FirstName:     u.FirstName(),
			LastName:      u.LastName(),
			Email:         u.Email(),
			Enabled:       u.Enabled(),
			CreatedAt:     u.CreatedAt(),
			LastUpdatedAt: u.LastUpdatedAt(),
		}))

	return err
}

func (r *UserSqlRepository) UpdateUser(u *user.User) error {
	columns := ColumnList{Users.FirstName, Users.LastName, Users.Email,
		Users.EmailVerified, Users.Enabled, Users.LastUpdatedAt}

	_, err := database.Execute(r.dbProvider, Users.UPDATE(columns).
		MODEL(model.Users{
			FirstName:     u.FirstName(),
			LastName:      u.LastName(),
			Email:         u.Email(),
			EmailVerified: u.EmailVerified(),
			Enabled:       u.Enabled(),
			LastUpdatedAt: u.LastUpdatedAt(),
		}).
		WHERE(Users.UserUUID.EQ(UUID(u.UUID()))))

	return err
}

func (r *UserSqlRepository) GetUserByEmail(email string) (*user.User, error) {
	return r.queryForUser(Users.Email.EQ(String(email)))
}

func (r *UserSqlRepository) GetUserByUUID(uuid uuid.UUID) (*user.User, error) {
	return r.queryForUser(Users.UserUUID.EQ(UUID(uuid)))
}

func (r *UserSqlRepository) queryForUser(predicate BoolExpression) (*user.User, error) {
	stmt := Users.
		LEFT_JOIN(UserRoles, UserRoles.UserID.EQ(Users.UserID)).
		LEFT_JOIN(Roles, Roles.RoleID.EQ(UserRoles.RoleID)).
		LEFT_JOIN(Permissions, Permissions.PermissionID.EQ(Roles.RoleID)).
		SELECT(Users.AllColumns, Roles.AllColumns, Permissions.AllColumns).
		WHERE(predicate)

	var dest []userQueryResult
	err := database.Query(r.dbProvider, stmt, &dest)
	if err != nil {
		return nil, err
	}

	if len(dest) > 0 {
		return mapUser(dest[0])
	}

	return nil, nil
}

type userQueryResult struct {
	model.Users

	Roles []*struct {
		model.Roles

		Permissions []*struct {
			model.Permissions
		}
	}
}

func mapUser(result userQueryResult) (*user.User, error) {
	roles, err := mapRoles(result)
	if err != nil {
		return nil, err
	}

	return user.FromDatabase(
		int(result.UserID),
		result.UserUUID,
		result.FirstName,
		result.LastName,
		result.Email,
		result.EmailVerified,
		result.Enabled,
		roles,
		result.CreatedAt,
		result.LastUpdatedAt,
	)
}

func mapRoles(dest userQueryResult) ([]*user.Role, error) {
	var roles []*user.Role
	var err error
	for _, r := range dest.Roles {
		if r == nil {
			break
		}
		var permissions []*user.Permission
		for _, p := range r.Permissions {
			if p == nil {
				break
			}
			permissions, err = appendPermission(permissions, p)
			if err != nil {
				return nil, err
			}
		}
		roles, err = appendRole(roles, r, permissions)
		if err != nil {
			return nil, err
		}
	}
	return roles, nil
}

func appendPermission(
	permissions []*user.Permission,
	p *struct{ model.Permissions },
) ([]*user.Permission, error) {
	var description string
	if p.Description != nil {
		description = *p.Description
	}
	permission, err := user.PermissionFromDatabase(
		int(p.PermissionID),
		p.PermissionUUID,
		p.Name,
		description,
		p.CreatedAt,
		p.LastUpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	permissions = append(permissions, permission)
	return permissions, nil
}

func appendRole(roles []*user.Role,
	r *struct {
		model.Roles
		Permissions []*struct{ model.Permissions }
	},
	permissions []*user.Permission,
) ([]*user.Role, error) {
	var description string
	if r.Description != nil {
		description = *r.Description
	}
	role, err := user.RoleFromDatabase(
		int(r.RoleID),
		r.RoleUUID,
		r.Name,
		description,
		permissions,
		r.CreatedAt,
		r.LastUpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	roles = append(roles, role)
	return roles, nil
}
