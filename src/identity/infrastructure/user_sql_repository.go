package infrastructure

import (
	"database/sql"
	"os"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/common/database"
	"github.com/jbenzshawel/go-sandbox/identity/domain"
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

func (r *UserSqlRepository) CreateUser(user domain.User, password string) (err error) {
	_, err = database.ExecuteInsert(r.dbProvider, Users.INSERT(Users.MutableColumns).
		MODEL(model.Users{
			UUID:          user.UUID,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Email:         user.Email,
			Password:      password,
			Enabled:       user.Enabled,
			CreatedAt:     user.CreatedAt,
			LastUpdatedAt: user.LastUpdatedAt,
		}))

	return err
}

func (r *UserSqlRepository) GetUserByEmail(email string) (*domain.User, error) {
	return r.queryForUser(
		SELECT(Users.AllColumns).
			FROM(Users).
			WHERE(Users.Email.EQ(String(email))),
	)
}

func (r *UserSqlRepository) GetUserByUUID(uuid uuid.UUID) (*domain.User, error) {
	return r.queryForUser(
		SELECT(Users.AllColumns).
			FROM(Users).
			WHERE(Users.UUID.EQ(UUID(uuid))),
	)
}

func (r *UserSqlRepository) queryForUser(stmt SelectStatement) (*domain.User, error) {
	var users []model.Users
	err := database.Query(r.dbProvider, stmt, &users)
	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		return mapToDomain(users[0]), nil
	}

	return nil, nil
}

func mapToDomain(user model.Users) *domain.User {
	return &domain.User{
		ID:            user.ID,
		UUID:          user.UUID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Enabled:       user.Enabled,
		CreatedAt:     user.CreatedAt,
		LastUpdatedAt: user.LastUpdatedAt,
	}
}
