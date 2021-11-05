package http

import (
	"backend/service/event/usecase"
	"testing"
)

func TestCreateEvent(t *testing.T) {
	useCaseMock := new(usecase.UseCaseMock)
	deliveryTest := NewDelivery(useCaseMock)

	userId := "1"

}
