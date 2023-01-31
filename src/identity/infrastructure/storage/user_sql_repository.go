package storage

import (
	"context"
	"database/sql"
	"os"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
	"github.com/jbenzshawel/go-sandbox/common/database"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user/permission"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user/role"
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

func (r *UserSqlRepository) Create(u *user.User) (err error) {
	db, err := r.dbProvider()
	if err != nil {
		return err
	}
	defer func() {
		closeErr := db.Close()
		err = cerror.CombineErrors(err, closeErr)
	}()
	txn, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := txn.Rollback()
			err = cerror.CombineErrors(err, rollbackErr)
		} else {
			commitErr := txn.Commit()
			err = cerror.CombineErrors(err, commitErr)
		}
	}()

	var createdUser model.Users
	stmt := Users.INSERT(Users.MutableColumns).
		MODEL(model.Users{
			UserUUID:      u.UUID(),
			FirstName:     u.FirstName(),
			LastName:      u.LastName(),
			Email:         u.Email(),
			Enabled:       u.Enabled(),
			CreatedAt:     u.CreatedAt(),
			LastUpdatedAt: u.LastUpdatedAt(),
		}).RETURNING(Users.AllColumns)
	err = stmt.QueryContext(context.Background(), txn, &createdUser)
	if err != nil {
		return err
	}

	for _, rl := range u.Roles() {
		stmt = UserRoles.INSERT(UserRoles.UserID, UserRoles.RoleID).
			MODEL(model.UserRoles{
				UserID: createdUser.UserID,
				RoleID: int32(rl.ID()),
			})
		_, err = stmt.ExecContext(context.Background(), txn)
		if err != nil {
			return err
		}
	}

	return
}

func (r *UserSqlRepository) Update(u *user.User) error {
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

func (r *UserSqlRepository) GetByEmail(email string) (*user.User, error) {
	return r.queryForUser(Users.Email.EQ(String(email)))
}

func (r *UserSqlRepository) GetByUUID(uuid uuid.UUID) (*user.User, error) {
	return r.queryForUser(Users.UserUUID.EQ(UUID(uuid)))
}

func (r *UserSqlRepository) queryForUser(predicate BoolExpression) (*user.User, error) {
	stmt := Users.
		LEFT_JOIN(UserRoles, UserRoles.UserID.EQ(Users.UserID)).
		LEFT_JOIN(Roles, Roles.RoleID.EQ(UserRoles.RoleID)).
		LEFT_JOIN(RolePermissions, RolePermissions.RoleID.EQ(Roles.RoleID)).
		LEFT_JOIN(Permissions, Permissions.PermissionID.EQ(RolePermissions.PermissionID)).
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

func mapRoles(dest userQueryResult) ([]*role.Role, error) {
	var roles []*role.Role
	var err error
	for _, r := range dest.Roles {
		if r == nil {
			break
		}
		var permissions []*permission.Permission
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
	permissions []*permission.Permission,
	p *struct{ model.Permissions },
) ([]*permission.Permission, error) {
	permit, err := permission.FromDatabase(
		int(p.PermissionID),
		p.Name,
	)
	if err != nil {
		return nil, err
	}
	permissions = append(permissions, permit)
	return permissions, nil
}

func appendRole(roles []*role.Role,
	r *struct {
		model.Roles
		Permissions []*struct{ model.Permissions }
	},
	permissions []*permission.Permission,
) ([]*role.Role, error) {
	rl, err := role.FromDatabase(
		int(r.RoleID),
		r.Name,
		permissions,
	)
	if err != nil {
		return nil, err
	}
	roles = append(roles, rl)
	return roles, nil
}
