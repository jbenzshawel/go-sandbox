package infrastructure

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/identity/domain"
	"github.com/jbenzshawel/go-sandbox/identity/infrastructure/gen/identity/identity/model"
	. "github.com/jbenzshawel/go-sandbox/identity/infrastructure/gen/identity/identity/table"
)

type DbConnectionFactory func() (*sql.DB, error)

type UserSqlRepository struct {
	dbConnection DbConnectionFactory
}

func NewUserSqlRepository(dbConnection func() (*sql.DB, error)) *UserSqlRepository {
	return &UserSqlRepository{
		dbConnection: dbConnection,
	}
}

func (r *UserSqlRepository) CreateUser(user domain.User, password string) (err error) {
	_, err = executeInsert(r.dbConnection, Users.INSERT(Users.MutableColumns).
		MODEL(model.Users{
			UUID:          user.UUID,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Email:         user.Email,
			Password:      password,
			Enabled:       true,
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
			WHERE(Users.UUID.EQ(String(uuid.String()))),
	)
}

func (r *UserSqlRepository) queryForUser(stmt SelectStatement) (*domain.User, error) {
	var users []model.Users
	err := query[*[]model.Users](r.dbConnection, stmt, &users)
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
		Enabled:       true,
		CreatedAt:     user.CreatedAt,
		LastUpdatedAt: user.LastUpdatedAt,
	}
}

// TODO: Move this to shared library/embedded struct?
func query[T any](dbConnection DbConnectionFactory, stmt SelectStatement, result T) (err error) {
	db, err := dbConnection()
	if err != nil {
		return err
	}
	defer func() {
		err = db.Close()
	}()

	err = stmt.Query(db, result)

	return err
}

func executeInsert(dbConnection DbConnectionFactory, stmt InsertStatement) (sql.Result, error) {
	db, err := dbConnection()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = db.Close()
	}()

	result, err := stmt.Exec(db)

	return result, err
}
