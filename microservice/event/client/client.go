package client

import (
	proto "backend/microservice/event/proto"
	"backend/pkg/models"
	"backend/service/event"
	"context"
)

const logMessage = "microservice:event:client:"

type EventService struct {
	repository event.Repository
}

func NewEventService(repository event.Repository) *EventService {
	return &EventService{
		repository: repository,
	}
}

func MakeProtoEvent(e *models.Event) *proto.Event {
	if e == nil {
		return &proto.Event{}
	}
	return &proto.Event{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Text:        e.Text,
		City:        e.City,
		Category:    e.Category,
		Viewed:      int32(e.Viewed),
		ImgUrl:      e.ImgUrl,
		Tag:         e.Tag,
		Date:        e.Date,
		Geo:         e.Geo,
		Address:     e.Address,
		AuthorId:    e.AuthorId,
		IsVisited:   e.IsVisited,
	}
}

func MakeProtoEvents(e []*models.Event) *proto.Events {
	if e == nil {
		return &proto.Events{}
	}
	result := make([]*proto.Event, len(e))
	for i, modelEvent := range e {
		result[i] = MakeProtoEvent(modelEvent)
	}
	return &proto.Events{
		Events: result,
	}
}

func MakeModelEvent(out *proto.Event) *models.Event {
	return &models.Event{
		ID:          out.ID,
		Title:       out.Title,
		Description: out.Description,
		Text:        out.Text,
		City:        out.City,
		Category:    out.Category,
		Viewed:      int(out.Viewed),
		ImgUrl:      out.ImgUrl,
		Tag:         out.Tag,
		Date:        out.Date,
		Geo:         out.Geo,
		Address:     out.Address,
		AuthorId:    out.AuthorId,
		IsVisited:   out.IsVisited,
	}
}

func MakeProtoInfo(info *models.Info) *proto.EmailInfo {
	if info == nil {
		return &proto.EmailInfo{}
	}
	return &proto.EmailInfo{
		Name: info.Name,
		Mail: info.Mail,
		Title: info.Title,
		ImgUrl: info.Img_url,
	}
}

func MakeProtoInfoArray(infoArray []*models.Info) *proto.EmailInfoArray {
	if infoArray == nil {
		return &proto.EmailInfoArray{}
	}
	result := make([]*proto.EmailInfo, len(infoArray))
	for i, info := range infoArray {
		result[i] = MakeProtoInfo(info)
	}
	return &proto.EmailInfoArray{
		InfoArray: result,
	}
}

func (c *EventService) CreateEvent(ctx context.Context, in *proto.Event) (*proto.EventId, error) {
	modelEvent := MakeModelEvent(in)
	eventId, err := c.repository.CreateEvent(modelEvent)
	out := &proto.EventId{
		ID: eventId,
	}
	return out, err
}

func (c *EventService) UpdateEvent(ctx context.Context, in *proto.UpdateEventRequest) (*proto.Empty, error) {
	modelEvent := MakeModelEvent(in.Event)
	userId := in.UserId
	err := c.repository.UpdateEvent(modelEvent, userId)
	out := &proto.Empty{}
	return out, err
}

func (c *EventService) DeleteEvent(ctx context.Context, in *proto.DeleteEventRequest) (*proto.Empty, error) {
	eventId := in.EventId
	userId := in.UserId
	err := c.repository.DeleteEvent(eventId, userId)
	out := &proto.Empty{}
	return out, err
}

func (c *EventService) GetEventById(ctx context.Context, in *proto.EventId) (*proto.Event, error) {
	eventId := in.ID
	modelEvent, err := c.repository.GetEventById(eventId)
	out := MakeProtoEvent(modelEvent)
	return out, err
}

func (c *EventService) GetEvents(ctx context.Context, in *proto.GetEventsRequest) (*proto.Events, error) {
	userId := in.UserId
	title := in.Title
	category := in.Category
	city := in.City
	date := in.Date
	tags := in.Tags
	modelEvents, err := c.repository.GetEvents(userId, title, category, city, date, tags)
	out := MakeProtoEvents(modelEvents)
	return out, err
}

func (c *EventService) GetVisitedEvents(ctx context.Context, in *proto.UserId) (*proto.Events, error) {
	userId := in.ID
	modelEvents, err := c.repository.GetVisitedEvents(userId)
	out := MakeProtoEvents(modelEvents)
	return out, err
}

func (c *EventService) GetCreatedEvents(ctx context.Context, in *proto.UserId) (*proto.Events, error) {
	userId := in.ID
	modelEvents, err := c.repository.GetCreatedEvents(userId)
	out := MakeProtoEvents(modelEvents)
	return out, err
}

func (c *EventService) Visit(ctx context.Context, in *proto.VisitRequest) (*proto.Empty, error) {
	eventId := in.EventId
	userId := in.UserId
	err := c.repository.Visit(eventId, userId)
	out := &proto.Empty{}
	return out, err
}

func (c *EventService) Unvisit(ctx context.Context, in *proto.VisitRequest) (*proto.Empty, error) {
	eventId := in.EventId
	userId := in.UserId
	err := c.repository.Unvisit(eventId, userId)
	out := &proto.Empty{}
	return out, err
}

func (c *EventService) IsVisited(ctx context.Context, in *proto.VisitRequest) (*proto.IsVisitedRequest, error) {
	eventId := in.EventId
	userId := in.UserId
	result, err := c.repository.IsVisited(eventId, userId)
	out := &proto.IsVisitedRequest{
		Result: result,
	}
	return out, err
}

func (c *EventService) GetCities(ctx context.Context, in *proto.Empty) (*proto.GetCitiesRequest, error) {
	result, err := c.repository.GetCities()
	out := &proto.GetCitiesRequest{
		Cities: result,
	}
	return out, err
}

func (c *EventService) EmailNotify(ctx context.Context, in *proto.EventId) (*proto.EmailInfoArray, error) {
	eventId := in.ID
	result, err := c.repository.EmailNotify(eventId)
	out := &proto.EmailInfoArray{
		InfoArray: MakeProtoInfoArray(result).InfoArray,
	}
	return out, err
}