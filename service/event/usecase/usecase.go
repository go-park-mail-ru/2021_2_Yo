package usecase

import (
	proto "backend/microservice/event/proto"
	log "backend/pkg/logger"
	"backend/pkg/models"
	error2 "backend/service/event/error"
	"context"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
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
		Address: 	 e.Address,
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
		Address: 	 out.Address,
		AuthorId:    out.AuthorId,
	}
}

func сityAndAddrByCoordinates(latitude, longitude string) (string, string)  {
	url := "https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address"
	url += "?lat="+latitude+"&lon="+longitude;
	log.Info(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("Accept","application/json")
	req.Header.Set("Authorization", "Token aaa00e3861df0b3fe38857306563ad4bee84550f")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	body, err := ioutil.ReadAll(resp.Body)  
	if err != nil {
		log.Error(err)
	}
	log.Debug(string(body))
	type Data struct {
		City string `json:"city,omittempty`
	}

	type AddrInfo struct {
		Value string `json:"value,omittempty`
		Unrestricted_value string `json:"unresticted_value,omitempty`
		Data Data `json:"data,omitempty"`
	}

	type Suggest struct {
		Suggestions []AddrInfo `json:"suggestions,omitempty"`
	}
	suggestions := Suggest{}
	
	err = json.Unmarshal(body, &suggestions)
	if err != nil {
		log.Error(err)
	}
	addr := suggestions.Suggestions[0].Value
	city := suggestions.Suggestions[0].Data.City
	
	return city, addr
}

func parseCoordinates(coords string) (string, string) {
	coordsArr := strings.Split(coords, " ")
	lat := coordsArr[0][1:len(coordsArr[0])-1]
	lng := coordsArr[1][:len(coordsArr[1])-1]
	log.Debug("x =", lat, "y =", lng)
	return lat, lng
}

func (a *UseCase) CreateEvent(e *models.Event) (string, error) {
	if e == nil || e.AuthorId == "" {
		return "", error2.ErrEmptyData
	}
	for i, tag := range e.Tag {
		e.Tag[i] = strings.ToLower(tag)
	}
	lat, lng := parseCoordinates(e.Geo)
	e.City,e.Address = сityAndAddrByCoordinates(lat,lng)
	
	in := MakeProtoEvent(e)
	log.Info("before repo")
	res, err := a.eventRepo.CreateEvent(context.Background(), in)
	log.Info(res)
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
	for i, tag := range e.Tag {
		e.Tag[i] = strings.ToLower(tag)
	}
	lat, lng := parseCoordinates(e.Geo)
	e.City,e.Address = сityAndAddrByCoordinates(lat,lng)

	in := &proto.UpdateEventRequest{
		Event:  MakeProtoEvent(e),
		UserId: userId,
	}
	_, err := a.eventRepo.UpdateEvent(context.Background(), in)
	log.Debug(logMessage+"UpdateEvent:HERE")
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

func (a *UseCase) GetCreatedEvents(userId string) ([]*models.Event, error) {
	if userId == "" {
		return nil, error2.ErrEmptyData
	}
	in := &proto.UserId{ID: userId}
	out, err := a.eventRepo.GetCreatedEvents(context.Background(), in)
	if err != nil {
		return nil, err
	}
	result := make([]*models.Event, len(out.Events))
	for i, protoEvent := range out.Events {
		result[i] = MakeModelEvent(protoEvent)
	}
	return result, nil
}
