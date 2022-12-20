package infrastructure

import (
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/identity/domain"
)

type UserMemoryRepository struct {
	users map[string]domain.User
	lock  *sync.RWMutex
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		users: map[string]domain.User{},
		lock:  &sync.RWMutex{},
	}
}

func (r *UserMemoryRepository) CreateUser(user domain.User, password string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, exists := r.users[user.Email]; exists {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	r.users[user.Email] = user

	return nil
}

func (r *UserMemoryRepository) GetUserByEmail(email string) (*domain.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	user, ok := r.users[email]
	if ok {
		return &user, nil
	}

	return nil, nil
}

func (r *UserMemoryRepository) GetUserByUUID(uuid uuid.UUID) (*domain.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, user := range r.users {
		if user.UUID == uuid {
			return &user, nil
		}
	}

	return nil, nil
}
