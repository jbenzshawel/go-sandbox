package storage

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user/role"
)

func TestUserRepository(t *testing.T) {
	repositories := createUserRepositories()
	for i := range repositories {
		r := repositories[i]

		t.Run(r.name, func(t *testing.T) {
			userRepo := r.repository
			var u *user.User
			var err error

			t.Run("Create", func(t *testing.T) {
				u, err = user.NewUser(
					"TestFirstName",
					"TestLastName",
					fmt.Sprintf("%s@test.com", uuid.New().String()),
					false,
					false,
				)
				require.NoError(t, err)
				require.NoError(t, u.SetUUID(uuid.New()))
				u.AddRole(role.Admin)
				err = userRepo.Create(u)
				require.NoError(t, err)
			})

			t.Run("GetByUUID", func(t *testing.T) {
				createdUser, err := userRepo.GetByUUID(u.UUID())
				require.NoError(t, err)
				require.NotNil(t, createdUser)
				assert.Greater(t, createdUser.ID(), 0)
				assertUserEqual(t, u, createdUser)
			})

			t.Run("GetByEmail", func(t *testing.T) {
				createdUser, err := userRepo.GetByEmail(u.Email())
				require.NoError(t, err)
				require.NotNil(t, createdUser)
				assert.Greater(t, createdUser.ID(), 0)
				assertUserEqual(t, u, createdUser)
			})

			if r.name == "PostgreSQL" {
				t.Run("User loads with permissions", func(t *testing.T) {
					createdUser, err := userRepo.GetByUUID(u.UUID())
					require.NoError(t, err)
					require.NotNil(t, createdUser)

					require.Len(t, createdUser.Roles(), 1)
					require.Len(t, createdUser.Roles()[0].Permissions(), 4)
				})
			}
		})
	}
}

func assertUserEqual(t *testing.T, expectedUser *user.User, actualUser *user.User) {
	assert.Equal(t, expectedUser.UUID(), actualUser.UUID())
	assert.Equal(t, expectedUser.FirstName(), actualUser.FirstName())
	assert.Equal(t, expectedUser.LastName(), actualUser.LastName())
	assert.Equal(t, expectedUser.Email(), actualUser.Email())
	assert.Equal(t, expectedUser.Enabled(), actualUser.Enabled())
	require.Len(t, actualUser.Roles(), 1)
	require.Len(t, expectedUser.Roles(), 1)
	assert.Equal(t, actualUser.Roles()[0].Type(), expectedUser.Roles()[0].Type())
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
