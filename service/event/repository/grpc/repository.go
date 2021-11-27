package grpc

import (
	eventGrpc "backend/microservice/event/proto"
	"backend/pkg/models"
	"context"
)

type Repository struct {
	client eventGrpc.EventServiceClient
}

func NewRepository(client eventGrpc.EventServiceClient) *Repository {
	return &Repository{
		client: client,
	}
}

func (s *Repository) CreateEvent(e *models.Event) (string, error) {
	in := &eventGrpc.Event{
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
	}
	out, err := s.client.CreateEvent(context.Background(), in)
	eventId := out.ID
	return eventId, err
}

func (s *Repository) UpdateEvent(e *models.Event, userId string) error {
	protoEvent := &eventGrpc.Event{
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
	}
	in := &eventGrpc.UpdateEventRequest{
		Event:  protoEvent,
		UserId: userId,
	}
	out, err := s.client.UpdateEvent(context.Background(), in)
	_ = out
	return err
}

func (s *Repository) DeleteEvent(eventId string, userId string) error {
	in := &eventGrpc.DeleteEventRequest{
		EventId: eventId,
		UserId:  userId,
	}
	out, err := s.client.DeleteEvent(context.Background(), in)
	_ = out
	return err
}

func (s *Repository) GetEventById(eventId string) (*models.Event, error) {
	in := &eventGrpc.EventId{
		ID: eventId,
	}
	out, err := s.client.GetEventById(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := &models.Event{
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
	}
	return result, err
}

func (s *Repository) GetEvents(title string, category string, city string, date string, tags []string) ([]*models.Event, error) {
	in := &eventGrpc.GetEventsRequest{
		Title:    title,
		Category: category,
		City:     city,
		Date:     date,
		Tags:     tags,
	}
	out, err := s.client.GetEvents(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := make([]*models.Event, len(out.Events))
	for i, protoEvent := range out.Events {
		result[i] = &models.Event{
			ID:          protoEvent.ID,
			Title:       protoEvent.Title,
			Description: protoEvent.Description,
			Text:        protoEvent.Text,
			City:        protoEvent.City,
			Category:    protoEvent.Category,
			Viewed:      int(protoEvent.Viewed),
			ImgUrl:      protoEvent.ImgUrl,
			Tag:         protoEvent.Tag,
			Date:        protoEvent.Date,
			Geo:         protoEvent.Geo,
			Address:     protoEvent.Address,
			AuthorId:    protoEvent.AuthorId,
		}
	}
	return result, err
}

func (s *Repository) GetVisitedEvents(userId string) ([]*models.Event, error) {
	in := &eventGrpc.UserId{
		ID: userId,
	}
	out, err := s.client.GetVisitedEvents(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := make([]*models.Event, len(out.Events))
	for i, protoEvent := range out.Events {
		result[i] = &models.Event{
			ID:          protoEvent.ID,
			Title:       protoEvent.Title,
			Description: protoEvent.Description,
			Text:        protoEvent.Text,
			City:        protoEvent.City,
			Category:    protoEvent.Category,
			Viewed:      int(protoEvent.Viewed),
			ImgUrl:      protoEvent.ImgUrl,
			Tag:         protoEvent.Tag,
			Date:        protoEvent.Date,
			Geo:         protoEvent.Geo,
			Address:     protoEvent.Address,
			AuthorId:    protoEvent.AuthorId,
		}
	}
	return result, err
}

func (s *Repository) GetCreatedEvents(authorId string) ([]*models.Event, error) {
	in := &eventGrpc.UserId{
		ID: authorId,
	}
	out, err := s.client.GetCreatedEvents(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := make([]*models.Event, len(out.Events))
	for i, protoEvent := range out.Events {
		result[i] = &models.Event{
			ID:          protoEvent.ID,
			Title:       protoEvent.Title,
			Description: protoEvent.Description,
			Text:        protoEvent.Text,
			City:        protoEvent.City,
			Category:    protoEvent.Category,
			Viewed:      int(protoEvent.Viewed),
			ImgUrl:      protoEvent.ImgUrl,
			Tag:         protoEvent.Tag,
			Date:        protoEvent.Date,
			Geo:         protoEvent.Geo,
			Address:     protoEvent.Address,
			AuthorId:    protoEvent.AuthorId,
		}
	}
	return result, err
}

func (s *Repository) Visit(eventId string, userId string) error {
	in := &eventGrpc.VisitRequest{
		EventId: eventId,
		UserId:  userId,
	}
	out, err := s.client.Visit(context.Background(), in)
	_ = out
	return err
}

func (s *Repository) Unvisit(eventId string, userId string) error {
	in := &eventGrpc.VisitRequest{
		EventId: eventId,
		UserId:  userId,
	}
	out, err := s.client.Unvisit(context.Background(), in)
	_ = out
	return err
}

func (s *Repository) IsVisited(eventId string, userId string) (bool, error) {
	in := &eventGrpc.VisitRequest{
		EventId: eventId,
		UserId:  userId,
	}
	out, err := s.client.IsVisited(context.Background(), in)
	result := out.Result
	return result, err
}

func (s *Repository) GetCities() ([]string, error) {
	in := &eventGrpc.Empty{}
	out, err := s.client.GetCities(context.Background(), in)
	result := out.Cities
	return result, err
}
