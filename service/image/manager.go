package image

import (
	"mime/multipart"
)

type Manager interface {
	SaveFile(userId string, fileName string, file multipart.File) error
}
