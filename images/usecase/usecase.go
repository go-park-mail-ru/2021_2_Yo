package usecase

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"backend/images"
	"github.com/satori/go.uuid"
)

type UseCase struct {
	imagesRepo images.Repository
}

func NewUseCase(imagesRepo images.Repository) *UseCase {
	return &UseCase{
		imagesRepo: imagesRepo,
	}
}

func (img *UseCase) SaveFile (userId string, fileName string,file multipart.File) {
	imgU := uuid.NewV4()
	s := strings.Split(fileName, ".")
	s[0] += imgU.String()
	newFileName := s[0] +s[1]
	dst, err := os.Create(filepath.Join("~/go/2021_2_Yo/static/images", filepath.Base(newFileName)))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		fmt.Println(err)
		return
	}
	img.imagesRepo.StoreImg(userId, fileName )
	
}