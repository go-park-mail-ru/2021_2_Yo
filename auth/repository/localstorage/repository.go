package localstorage

import (
	"backend/auth"
	"backend/models"
	"strconv"
	"sync"
)

type RepositoryUserLocalStorage struct {
	users []*User
	mutex *sync.Mutex
}

func NewRepositoryUserLocalStorage() *RepositoryUserLocalStorage {
	result := &RepositoryUserLocalStorage{
		users: make([]*User, 0),
		mutex: new(sync.Mutex),
	}
	return result
}

func (s *RepositoryUserLocalStorage) CreateUser(user *models.User) error {
	newUser := toLocalstorageUser(user)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, tempUser := range s.users {
		if tempUser.Mail == user.Mail {
			return auth.ErrUserExists
		}
	}
	if len(s.users) > 0 {
		newUser.ID = s.users[len(s.users)-1].ID + 1
	} else {
		newUser.ID = 0
	}
	s.users = append(s.users, newUser)
	return nil
}

func (s *RepositoryUserLocalStorage) GetUser(mail, password string) (*models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, user := range s.users {
		if user.Mail == mail && user.Password == password {
			return toModelUser(user), nil
		}
	}
	return nil, auth.ErrUserNotFound
}

func (s *RepositoryUserLocalStorage) GetUserById(userId string) (*models.User, error) {
	userIdInt, _ := strconv.Atoi(userId)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, user := range s.users {
		if user.ID == userIdInt {
			return toModelUser(user), nil
		}
	}
	return nil, auth.ErrUserNotFound
}
