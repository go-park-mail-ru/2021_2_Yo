package images

import "mime/multipart"

type UseCase interface {
	SaveFile (userId string, fileName string, file multipart.File)
}