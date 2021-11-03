package image

import (
	"github.com/stretchr/testify/mock"
	"mime/multipart"
)

type ManagerMock struct {
	mock.Mock
}

func (m *ManagerMock) SaveFile(userId string, fileName string, file multipart.File) error {
	args := m.Called(userId, fileName, file)
	return args.Error(0)
}
