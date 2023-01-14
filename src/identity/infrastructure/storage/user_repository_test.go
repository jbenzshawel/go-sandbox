package storage

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
)

func TestUserRepository(t *testing.T) {
	repositories := createUserRepositories()
	for i := range repositories {
		r := repositories[i]

		t.Run(r.name, func(t *testing.T) {
			userRepo := r.repository
			var u *user.User
			var err error

			t.Run("InsertUser", func(t *testing.T) {
				u, err = user.NewUser(
					"TestFirstName",
					"TestLastName",
					fmt.Sprintf("%s@test.com", uuid.New().String()),
					false,
					false,
				)
				require.NoError(t, err)
				require.NoError(t, u.SetUUID(uuid.New()))
				err = userRepo.InsertUser(u)
				require.NoError(t, err)
			})

			t.Run("GetUserByUUID", func(t *testing.T) {
				createdUser, err := userRepo.GetUserByUUID(u.UUID())
				require.NoError(t, err)
				require.NotNil(t, createdUser)
				assert.Greater(t, createdUser.ID(), int32(0))
				assertUserEqual(t, u, createdUser)
			})

			t.Run("GetUserByEmail", func(t *testing.T) {
				createdUser, err := userRepo.GetUserByEmail(u.Email())
				require.NoError(t, err)
				require.NotNil(t, createdUser)
				assert.Greater(t, createdUser.ID(), int32(0))
				assertUserEqual(t, u, createdUser)
			})
		})
	}
}

func assertUserEqual(t *testing.T, expectedUser *user.User, actualUser *user.User) {
	assert.Equal(t, expectedUser.UUID(), actualUser.UUID())
	assert.Equal(t, expectedUser.FirstName(), actualUser.FirstName())
	assert.Equal(t, expectedUser.LastName(), actualUser.LastName())
	assert.Equal(t, expectedUser.Email(), actualUser.Email())
	assert.Equal(t, expectedUser.Enabled(), actualUser.Enabled())
	assert.WithinDuration(t, expectedUser.CreatedAt(), actualUser.CreatedAt(), 0)
	assert.WithinDuration(t, expectedUser.LastUpdatedAt(), actualUser.LastUpdatedAt(), 0)
}

type repository struct {
	name       string
	repository user.Repository
}

func createUserRepositories() []repository {
	repositories := []repository{
		{
			name:       "Memory",
			repository: NewUserMemoryRepository(),
		},
	}
	if userSqlRepo, ok := TryCreateUserSqlRepository(); ok {
		repositories = append(repositories, repository{
			name:       "PostgreSQL",
			repository: userSqlRepo,
		})
	}
	return repositories
}
