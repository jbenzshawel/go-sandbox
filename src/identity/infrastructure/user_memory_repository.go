package infrastructure

import (
	"fmt"
	"sync"

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

func (r *UserMemoryRepository) RegisterUser(user domain.User, password string) error {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if _, exists := r.users[user.Email]; exists {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	r.users[user.Email] = user

	return nil
}
