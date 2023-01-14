package storage

import (
	"sync"

	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
)

type UserMemoryRepository struct {
	users  map[string]user.User
	lock   *sync.RWMutex
	lastId int32
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		users:  map[string]user.User{},
		lock:   &sync.RWMutex{},
		lastId: 0,
	}
}

func (r *UserMemoryRepository) InsertUser(user *user.User) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.lastId++
	err := user.SetID(r.lastId)
	if err != nil {
		return err
	}
	r.users[user.Email()] = *user

	return nil
}

func (r *UserMemoryRepository) UpdateUser(user *user.User) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.users[user.Email()] = *user

	return nil
}

func (r *UserMemoryRepository) GetUserByEmail(email string) (*user.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	u, ok := r.users[email]
	if ok {
		return &u, nil
	}

	return nil, nil
}

func (r *UserMemoryRepository) GetUserByUUID(uuid uuid.UUID) (*user.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, u := range r.users {
		if u.UUID() == uuid {
			return &u, nil
		}
	}

	return nil, nil
}
