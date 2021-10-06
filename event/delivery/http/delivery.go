package http

import (
	"backend/event"
	"backend/response"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Delivery struct {
	useCase event.UseCase
}

func NewDelivery(useCase event.UseCase) *Delivery {
	return &Delivery{
		useCase: useCase,
	}
}

func (h *Delivery) List(w http.ResponseWriter, r *http.Request) {
	eventsList, err := h.useCase.List()
	if err != nil {
		log.Error("List : got error", err)
		response.SendResponse(w, response.ErrorResponse("Can't get list of events"))
		return
	}
	response.SendResponse(w, response.EventsListResponse(eventsList))
}
