package images

import (
	"backend/images/repository"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	 
	"github.com/satori/go.uuid"
)

type Manager struct {
	imagesRepo repository.Repository
}

func NewManager(imagesRepo repository.Repository) *Manager {
	return &Manager{
		imagesRepo: imagesRepo,
	}
}

func (img *Manager) SaveFile (userId string, fileName string,file multipart.File) error {
	imgU := uuid.NewV4()
	s := strings.Split(fileName, ".")
	s[0] += imgU.String()
	newFileName := s[0] + "." + s[1]
	dst, err := os.Create(filepath.Join("/home/ubuntu/go/2021_2_Yo/static/images", filepath.Base(newFileName)))
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		fmt.Println(err)
		return err
	}
	err = img.imagesRepo.StoreImage(userId,fileName)
	if err != nil {
		return err
	}
	return nil
}