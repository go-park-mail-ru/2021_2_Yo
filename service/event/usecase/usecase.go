package usecase

import (
	proto "backend/microservice/event/proto"
	"backend/pkg/models"
	error2 "backend/service/event/error"
	"context"
	"strings"
)

const logMessage = "service:event:usecase:"

type UseCase struct {
	eventRepo proto.RepositoryClient
}

func NewUseCase(eventRepo proto.RepositoryClient) *UseCase {
	return &UseCase{
		eventRepo: eventRepo,
	}
}

func MakeProtoEvent(e *models.Event) *proto.Event {
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
		AuthorId:    e.AuthorId,
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
		AuthorId:    out.AuthorId,
	}
}

func (a *UseCase) CreateEvent(e *models.Event) (string, error) {
	if e == nil || e.AuthorId == "" {
		return "", error2.ErrEmptyData
	}
	for _, tag := range e.Tag {
		tag = strings.ToLower(tag)
	}
	in := MakeProtoEvent(e)
	res, err := a.eventRepo.CreateEvent(context.Background(), in)
	if err != nil {
		return "", err
	}
	return res.ID, nil
}

func (a *UseCase) UpdateEvent(e *models.Event, userId string) error {
	if e == nil || userId == "" {
		return error2.ErrEmptyData
	}
	if e.ID == "" {
		return error2.ErrEmptyData
	}
	for _, tag := range e.Tag {
		tag = strings.ToLower(tag)
	}
	in := &proto.UpdateEventRequest{
		Event:  MakeProtoEvent(e),
		UserId: userId,
	}
	_, err := a.eventRepo.UpdateEvent(context.Background(), in)
	return err
}

func (a *UseCase) DeleteEvent(userId string, eventID string) error {
	if userId == "" || eventID == "" {
		return error2.ErrEmptyData
	}
	in := &proto.DeleteEventRequest{
		EventId: eventID,
		UserId:  userId,
	}
	_, err := a.eventRepo.DeleteEvent(context.Background(), in)
	return err
}

func (a *UseCase) GetEventById(eventId string) (*models.Event, error) {
	if eventId == "" {
		return nil, error2.ErrEmptyData
	}
	in := &proto.EventId{ID: eventId}
	out, err := a.eventRepo.GetEventById(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := MakeModelEvent(out)
	return result, nil
}

func (a *UseCase) GetEvents(title string, category string, tags []string) ([]*models.Event, error) {
	if tags != nil && tags[0] == "" {
		tags = nil
	}
	in := &proto.GetEventsRequest{
		Title:    title,
		Category: category,
		Tags:     tags,
	}
	out, err := a.eventRepo.GetEvents(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := make([]*models.Event, len(out.Events))
	for i, protoEvent := range out.Events {
		result[i] = MakeModelEvent(protoEvent)
	}
	return result, nil
}

func (a *UseCase) GetEventsFromAuthor(authorId string) ([]*models.Event, error) {
	if authorId == "" {
		return nil, error2.ErrEmptyData
	}
	in := &proto.AuthorId{ID: authorId}
	out, err := a.eventRepo.GetEventsFromAuthor(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := make([]*models.Event, len(out.Events))
	for i, protoEvent := range out.Events {
		result[i] = MakeModelEvent(protoEvent)
	}
	return result, nil
}

func (a *UseCase) Visit(eventId string, userId string) error {
	if eventId == "" || userId == "" {
		return error2.ErrEmptyData
	}
	in := &proto.VisitRequest{
		EventId: eventId,
		UserId:  userId,
	}
	_, err := a.eventRepo.Visit(context.Background(), in)
	return err
}

func (a *UseCase) GetVisitedEvents(userId string) ([]*models.Event, error) {
	if userId == "" {
		return nil, error2.ErrEmptyData
	}
	in := &proto.UserId{
		ID: userId,
	}
	out, err := a.eventRepo.GetVisitedEvents(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := make([]*models.Event, len(out.Events))
	for i, protoEvent := range out.Events {
		result[i] = MakeModelEvent(protoEvent)
	}
	return result, nil
}
