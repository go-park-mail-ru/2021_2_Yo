package images

import (
	"backend/images/models"
	"net/http"
	"fmt"
	log "backend/logger"
	"backend/images"
)

var ImgPaths []models.Imgpath

type Delivery struct {
	useCase images.UseCase
}

func NewDelivery(useCase images.UseCase) *Delivery {
	return &Delivery{
		useCase: useCase,
	}
}

func (img *Delivery) uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(1 << 2)
	userId := r.Context().Value("userId").(string)
	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()

	log.Info("Uploaded File: %+v\n", handler.Filename)
	log.Info("File Size: %+v\n", handler.Size)
	log.Info("MIME Header: %+v\n", handler.Header)

	img.useCase.SaveFile(userId, handler.Filename, file)

	log.Info(w, "Successfully Uploaded File\n"+"")
}