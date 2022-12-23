package infrastructure

import (
	"sync"

	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/identity/domain"
)

type UserMemoryRepository struct {
	users  map[string]domain.User
	lock   *sync.RWMutex
	lastId int32
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		users:  map[string]domain.User{},
		lock:   &sync.RWMutex{},
		lastId: 0,
	}
}

func (r *UserMemoryRepository) CreateUser(user domain.User, password string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.lastId++
	user.ID = r.lastId
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
