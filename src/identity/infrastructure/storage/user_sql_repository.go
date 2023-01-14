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

func (r *UserSqlRepository) InsertUser(u *user.User) error {
	_, err := database.ExecuteInsert(r.dbProvider, Users.INSERT(Users.MutableColumns).
		MODEL(model.Users{
			UUID:          u.UUID(),
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

	_, err := database.ExecuteUpdate(r.dbProvider, Users.UPDATE(columns).
		MODEL(model.Users{
			FirstName:     u.FirstName(),
			LastName:      u.LastName(),
			Email:         u.Email(),
			EmailVerified: u.EmailVerified(),
			Enabled:       u.Enabled(),
			LastUpdatedAt: u.LastUpdatedAt(),
		}).
		WHERE(Users.UUID.EQ(UUID(u.UUID()))))

	return err
}

func (r *UserSqlRepository) GetUserByEmail(email string) (*user.User, error) {
	return r.queryForUser(
		SELECT(Users.AllColumns).
			FROM(Users).
			WHERE(Users.Email.EQ(String(email))),
	)
}

func (r *UserSqlRepository) GetUserByUUID(uuid uuid.UUID) (*user.User, error) {
	return r.queryForUser(
		SELECT(Users.AllColumns).
			FROM(Users).
			WHERE(Users.UUID.EQ(UUID(uuid))),
	)
}

func (r *UserSqlRepository) queryForUser(stmt SelectStatement) (*user.User, error) {
	var users []model.Users
	err := database.Query(r.dbProvider, stmt, &users)
	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		return mapToDomain(users[0])
	}

	return nil, nil
}

func mapToDomain(u model.Users) (*user.User, error) {
	return user.FromDatabase(
		u.ID,
		u.UUID,
		u.FirstName,
		u.LastName,
		u.Email,
		u.EmailVerified,
		u.Enabled,
		u.CreatedAt,
		u.LastUpdatedAt,
	)
}
