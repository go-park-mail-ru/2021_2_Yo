package usecase

import (
	log "backend/pkg/logger"
	"backend/pkg/models"
	"backend/service/event"
	error2 "backend/service/event/error"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const logMessage = "service:event:usecase:"

type UseCase struct {
	repository event.Repository
}

func NewUseCase(repository event.Repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}

func cityAndAddrByCoordinates(latitude, longitude string) (string, string, error) {
	url := "https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address"
	url += "?lat=" + latitude + "&lon=" + longitude
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token aaa00e3861df0b3fe38857306563ad4bee84550f")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	type Data struct {
		City string `json:"city,omittempty`
	}

	type AddrInfo struct {
		Value              string `json:"value,omittempty`
		Unrestricted_value string `json:"unresticted_value,omitempty`
		Data               Data   `json:"data,omitempty"`
	}

	type Suggest struct {
		Suggestions []AddrInfo `json:"suggestions,omitempty"`
	}
	suggestions := Suggest{}

	err = json.Unmarshal(body, &suggestions)
	if err != nil {
		return "", "", err
	}

	if len(suggestions.Suggestions) == 0 {
		return "", "", errors.New("can't get city and address from coordinates")
	}

	addr := suggestions.Suggestions[0].Value
	city := suggestions.Suggestions[0].Data.City

	return city, addr, nil
}

func parseCoordinates(coords string) (string, string) {
	coordsArr := strings.Split(coords, " ")
	lat := coordsArr[0][1 : len(coordsArr[0])-1]
	lng := coordsArr[1][:len(coordsArr[1])-1]
	return lat, lng
}

func (a *UseCase) CreateEvent(e *models.Event) (string, error) {
	if e == nil || e.AuthorId == "" {
		return "", error2.ErrEmptyData
	}
	lat, lng := parseCoordinates(e.Geo)
	city, address, err := cityAndAddrByCoordinates(lat, lng)
	if err != nil {
		log.Error(logMessage+"CreateEvent:err = ", err)
	} else {
		e.City = city
		e.Address = address
	}
	for i, tag := range e.Tag {
		e.Tag[i] = strings.ToLower(tag)
	}
	log.Debug(logMessage+"CreateEvent:e.ImgUrl = ", e.ImgUrl)
	return a.repository.CreateEvent(e)
}

func (a *UseCase) UpdateEvent(e *models.Event, userId string) error {
	if e == nil || userId == "" || e.ID == "" {
		return error2.ErrEmptyData
	}
	lat, lng := parseCoordinates(e.Geo)
	city, address, err := cityAndAddrByCoordinates(lat, lng)
	if err != nil {
		log.Error(logMessage+"CreateEvent:err = ", err)
	} else {
		e.City = city
		e.Address = address
	}
	for i, tag := range e.Tag {
		e.Tag[i] = strings.ToLower(tag)
	}
	return a.repository.UpdateEvent(e, userId)
}

func (a *UseCase) DeleteEvent(eventID string, userId string) error {
	if userId == "" || eventID == "" {
		return error2.ErrEmptyData
	}
	return a.repository.DeleteEvent(eventID, userId)
}

func (a *UseCase) GetEventById(eventId string) (*models.Event, error) {
	if eventId == "" {
		return nil, error2.ErrEmptyData
	}
	return a.repository.GetEventById(eventId)
}

func (a *UseCase) GetEvents(title string, category string, city string, date string, tags []string) ([]*models.Event, error) {
	if tags != nil && tags[0] == "" {
		tags = nil
	}
	return a.repository.GetEvents(title, category, city, date, tags)
}

func (a *UseCase) GetVisitedEvents(userId string) ([]*models.Event, error) {
	if userId == "" {
		return nil, error2.ErrEmptyData
	}
	return a.repository.GetVisitedEvents(userId)
}

func (a *UseCase) GetCreatedEvents(userId string) ([]*models.Event, error) {
	if userId == "" {
		return nil, error2.ErrEmptyData
	}
	return a.repository.GetCreatedEvents(userId)
}

func (a *UseCase) Visit(eventId string, userId string) error {
	if eventId == "" || userId == "" {
		return error2.ErrEmptyData
	}
	return a.repository.Visit(eventId, userId)
}

func (a *UseCase) Unvisit(eventId string, userId string) error {
	if eventId == "" || userId == "" {
		return error2.ErrEmptyData
	}
	return a.repository.Unvisit(eventId, userId)
}

func (a *UseCase) IsVisited(eventId string, userId string) (bool, error) {
	if eventId == "" || userId == "" {
		return false, error2.ErrEmptyData
	}
	return a.repository.IsVisited(eventId, userId)
}

func (a *UseCase) GetCities() ([]string, error) {
	return a.repository.GetCities()
}
