package localstorage

import (
	"backend/models"
	"errors"
	"sync"
)

type RepositoryUserLocalStorage struct {
	users []*User
	mutex *sync.Mutex
}

//Функции, создающие все хендлеры, юзкейсы и репозитории, вызываются непосредственно перед инициализированием сервера
func NewRepositoryUserLocalStorage() *RepositoryUserLocalStorage {
	return &RepositoryUserLocalStorage{
		users: make([]*User, 0),
		mutex: new(sync.Mutex),
	}
}

func (s *RepositoryUserLocalStorage) CreateUser(user *models.User) error {
	newUser := toLocalstorageUser(user)
	s.mutex.Lock()
	if len(s.users) > 0 {
		newUser.ID = s.users[len(s.users)-1].ID + 1
	} else {
		newUser.ID = 0
	}
	s.users = append(s.users, newUser)
	s.mutex.Unlock()
	return nil
}

func (s *RepositoryUserLocalStorage) GetUser(username, password string) (*models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, user := range s.users {
		if user.Username == username && user.Password == password {
			return toModelUser(user), nil
		}
	}

	return nil, errors.New("No user found")
}
