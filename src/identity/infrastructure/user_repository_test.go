package infrastructure

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
	t.Parallel()

	repositories := createUserRepositories()
	for i := range repositories {
		r := repositories[i]

		t.Run(r.name, func(t *testing.T) {
			t.Parallel()

			t.Run("testCreateUser", func(t *testing.T) {
				t.Parallel()
				testCreateUser(t, r.repository)
			})
		})
	}
}

func testCreateUser(t *testing.T, userRepository domain.UserRepository) {
	t.Helper()

	user := domain.User{
		UUID:          uuid.New(),
		FirstName:     "TestFirstName",
		LastName:      "TestLastName",
		Email:         fmt.Sprintf("%s.com", uuid.New().String()),
		Enabled:       false,
		CreatedAt:     time.Now().In(&time.Location{}),
		LastUpdatedAt: time.Now().In(&time.Location{}),
	}
	err := userRepository.CreateUser(user, uuid.New().String())
	require.NoError(t, err)

	createdUser, err := userRepository.GetUserByUUID(user.UUID)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	assert.Greater(t, createdUser.ID, int32(0))
	assert.Equal(t, user.UUID, createdUser.UUID)
	assert.Equal(t, user.FirstName, createdUser.FirstName)
	assert.Equal(t, user.LastName, createdUser.LastName)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Enabled, createdUser.Enabled)
	assert.WithinDuration(t, user.CreatedAt, createdUser.CreatedAt, 0)
	assert.WithinDuration(t, user.LastUpdatedAt, createdUser.LastUpdatedAt, 0)
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
