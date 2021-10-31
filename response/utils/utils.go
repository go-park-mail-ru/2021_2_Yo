package utils

import (
	log "backend/logger"
	"backend/response"
	"net/http"
)

func CheckIfNoError(w *http.ResponseWriter, err error, msg string, status response.HttpStatus) bool {
	if err != nil {
		log.Error(msg+"err =", err)
		response.ErrorResponse("")
		//response.SendResponse(*w, response.StatusResponse(status))
		return false
	}
	return true
}
