package response

import (
	log "backend/pkg/logger"
	models2 "backend/pkg/models"
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/go-sanitize/sanitize"
	"io"
	"net/http"
)

var (
	ErrJSONDecoding   = errors.New("data decoding error")
	ErrValidation     = errors.New("data validation error")
	ErrSanitizing     = errors.New("data sanitizing error")
	ErrSanitizerError = errors.New("internal sanitizing package error")
)

func ValidateAndSanitize(object interface{}) error {
	s, err := sanitize.New()
	if err != nil {
		return ErrSanitizerError
	}
	err = s.Sanitize(object)
	if err != nil {
		return ErrSanitizing
	}
	valid, err := govalidator.ValidateStruct(object)
	if err != nil || !valid {
		return ErrValidation
	}
	return nil
}

func GetEventFromRequest(r io.Reader) (*models2.Event, error) {
	eventInput := new(models2.EventResponseBody)
	err := json.NewDecoder(r).Decode(eventInput)
	if err != nil {
		return nil, ErrJSONDecoding
	}
	err = ValidateAndSanitize(eventInput)
	if err != nil {
		return nil, err
	}
	result := &models2.Event{
		ID:          eventInput.ID,
		Title:       eventInput.Title,
		Description: eventInput.Description,
		Text:        eventInput.Text,
		City:        eventInput.City,
		Category:    eventInput.Category,
		Viewed:      eventInput.Viewed,
		ImgUrl:      eventInput.ImgUrl,
		Tag:         eventInput.Tag,
		Date:        eventInput.Date,
		Geo:         eventInput.Geo,
	}
	return result, nil
}

func GetUserFromRequest(r io.Reader) (*models2.User, error) {
	message := logMessage + "GetUserFromRequest:"
	_ = message
	userInput := new(models2.UserResponseBody)
	err := json.NewDecoder(r).Decode(userInput)
	if err != nil {
		return nil, ErrJSONDecoding
	}
	err = ValidateAndSanitize(userInput)
	if err != nil {
		return nil, err
	}
	result := &models2.User{
		Name:     userInput.Name,
		Surname:  userInput.Surname,
		Mail:     userInput.Mail,
		Password: userInput.Password,
		About:    userInput.About,
	}
	valid, err := govalidator.ValidateStruct(result)
	if err != nil || !valid {
		return nil, ErrValidation
	}
	return result, nil
}

func MakeEventResponseBody(e *models2.Event) models2.EventResponseBody {
	return models2.EventResponseBody{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Text:        e.Text,
		City:        e.City,
		Category:    e.Category,
		Viewed:      e.Viewed,
		ImgUrl:      e.ImgUrl,
		Tag:         e.Tag,
		Date:        e.Date,
		Geo:         e.Geo,
		AuthorID:    e.AuthorId,
	}
}

func MakeEventListResponseBody(events []*models2.Event) models2.EventListResponseBody {
	result := make([]models2.EventResponseBody, len(events))
	for i := 0; i < len(events); i++ {
		result[i] = MakeEventResponseBody(events[i])
	}
	return models2.EventListResponseBody{
		Events: result,
	}
}

func MakeUserResponseBody(u *models2.User) models2.UserResponseBody {
	return models2.UserResponseBody{
		ID:       u.ID,
		Name:     u.Name,
		Surname:  u.Surname,
		About:    u.About,
		ImgUrl:   u.ImgUrl,
		Mail:     u.Mail,
		Password: u.Password,
	}
}

func MakeUserListResponseBody(users []*models2.User) models2.UserListResponseBody {
	result := make([]models2.UserResponseBody, len(users))
	for i := 0; i < len(users); i++ {
		result[i] = MakeUserResponseBody(users[i])
	}
	return models2.UserListResponseBody{
		Users: result,
	}
}

func SendResponse(w http.ResponseWriter, response interface{}) {
	message := logMessage + "SendResponse:"
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(response)
	if err != nil {
		log.Error(message+"err =", err)
		return
	}
	w.Write(b)
}
