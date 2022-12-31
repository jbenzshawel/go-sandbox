package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jbenzshawel/go-sandbox/identity/domain"
)

func TestUserRepository(t *testing.T) {
	repositories := createUserRepositories()
	for i := range repositories {
		r := repositories[i]

		t.Run(r.name, func(t *testing.T) {
			userRepo := r.repository
			var user domain.User

			t.Run("InsertUser", func(t *testing.T) {
				user = domain.User{
					UUID:          uuid.New(),
					FirstName:     "TestFirstName",
					LastName:      "TestLastName",
					Email:         fmt.Sprintf("%s@test.com", uuid.New().String()),
					Enabled:       false,
					CreatedAt:     time.Now().In(&time.Location{}),
					LastUpdatedAt: time.Now().In(&time.Location{}),
				}
				err := userRepo.InsertUser(user)
				require.NoError(t, err)
			})

			t.Run("GetUserByUUID", func(t *testing.T) {
				createdUser, err := userRepo.GetUserByUUID(user.UUID)
				require.NoError(t, err)
				require.NotNil(t, createdUser)
				assert.Greater(t, createdUser.ID, int32(0))
				assertUserEqual(t, &user, createdUser)
			})

			t.Run("GetUserByEmail", func(t *testing.T) {
				createdUser, err := userRepo.GetUserByEmail(user.Email)
				require.NoError(t, err)
				require.NotNil(t, createdUser)
				assert.Greater(t, createdUser.ID, int32(0))
				assertUserEqual(t, &user, createdUser)
			})
		})
	}
}

func assertUserEqual(t *testing.T, expectedUser *domain.User, actualUser *domain.User) {
	assert.Equal(t, expectedUser.UUID, actualUser.UUID)
	assert.Equal(t, expectedUser.FirstName, actualUser.FirstName)
	assert.Equal(t, expectedUser.LastName, actualUser.LastName)
	assert.Equal(t, expectedUser.Email, actualUser.Email)
	assert.Equal(t, expectedUser.Enabled, actualUser.Enabled)
	assert.WithinDuration(t, expectedUser.CreatedAt, actualUser.CreatedAt, 0)
	assert.WithinDuration(t, expectedUser.LastUpdatedAt, actualUser.LastUpdatedAt, 0)
}

type repository struct {
	name       string
	repository domain.UserRepository
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
