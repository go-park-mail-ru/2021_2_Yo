package http

import (
	"backend/eventsManager"
	"backend/response"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type HandlerEventsManager struct {
	useCase eventsManager.UseCaseEventsManager
}

func NewHandlerEventsManager(useCase eventsManager.UseCaseEventsManager) *HandlerEventsManager {
	return &HandlerEventsManager{
		useCase: useCase,
	}
}

func (h *HandlerEventsManager) List(w http.ResponseWriter, r *http.Request) {
	eventsList, err := h.useCase.List()
	if err != nil {
		log.Error("List : got error", err)
		response.SendResponse(w, response.ErrorResponse("Can't get list of events"))
		return
	}
	response.SendResponse(w, response.EventsListResponse(eventsList))
}
