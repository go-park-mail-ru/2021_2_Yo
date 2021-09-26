package localstorage

import (
	"backend/auth"
	"backend/models"
	"sync"
)

type RepositoryUserLocalStorage struct {
	users []*User
	mutex *sync.Mutex
}

//Функции, создающие все хендлеры, юзкейсы и репозитории, вызываются непосредственно перед инициализированием сервера
func NewRepositoryUserLocalStorage() *RepositoryUserLocalStorage {
	result := &RepositoryUserLocalStorage{
		users: make([]*User, 1),
		mutex: new(sync.Mutex),
	}
	result.users[0] = &User{0,"Dasha","Petrova","funnyduck@yandex.ru","1234567890"}
	return result
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

func (s *RepositoryUserLocalStorage) GetUser(mail,password string) (*models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, user := range s.users {
		if user.Mail == mail && user.Password == password {
			return toModelUser(user), nil
		}
	}

	return nil, auth.ErrUserNotFound
}

func (s *RepositoryUserLocalStorage) List() []*models.User {
	var UsersFromStorage []*models.User
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, user := range s.users {
		UsersFromStorage = append(UsersFromStorage, toModelUser(user))
	}
	return UsersFromStorage
}
