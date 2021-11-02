package csrf

import "github.com/stretchr/testify/mock"

type ManagerMock struct {
	mock.Mock
}

func (m *ManagerMock) Create(userId string) (string, error) {
	args := m.Called(userId)
	return "", args.Error(1)
}

func (m *ManagerMock) Check(sessionId string) (string, error) {
	args := m.Called(sessionId)
	return "", args.Error(1)
}

func (m *ManagerMock) Delete(sessionId string) error {
	args := m.Called(sessionId)
	return args.Error(0)
}
