package userrepo

import (
	"errors"
	"homework10/internal/entities"
	"homework10/internal/util"
	"sync"
)

var ErrEmptyUser = errors.New("user is empty")

type UserRepository interface {
	AddUser(user entities.User) (int64, error)
	EditUser(setUser entities.User) (*entities.User, error)
	GetUserByID(id int64) (*entities.User, error)
	DeleteUser(id int64) error
}

type mapRepository struct {
	rep    map[int64]entities.User
	mutex  sync.Mutex
	rMutex sync.RWMutex
	util.UID
}

func (m *mapRepository) AddUser(user entities.User) (int64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	const notValidID = -1
	id, err := m.UID.GenerateID()
	if err != nil {
		return notValidID, err
	}

	user.ID = id
	m.rep[id] = user
	return user.ID, nil
}

func (m *mapRepository) EditUser(setUser entities.User) (*entities.User, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, err := m.GetUserByID(setUser.ID); err != nil {
		return &setUser, err
	}

	m.rep[setUser.ID] = setUser
	return &setUser, nil
}

// GetUserByID получает пользователя по ID. Реализована проверка на существование пользователя
func (m *mapRepository) GetUserByID(id int64) (*entities.User, error) {
	m.rMutex.RLock()
	defer m.rMutex.RUnlock()

	empty := &entities.User{}
	user := m.rep[id]
	if user == (*empty) {
		return &user, ErrEmptyUser
	}
	return &user, nil
}

func (m *mapRepository) DeleteUser(id int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, err := m.GetUserByID(id); err != nil {
		return err
	}
	delete(m.rep, id)
	return nil
}
func New() UserRepository {
	return &mapRepository{
		rep: make(map[int64]entities.User),
		UID: util.UID{Id: -1}}
}
