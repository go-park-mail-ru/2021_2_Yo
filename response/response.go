package response

import (
	log "backend/logger"
	"backend/models"
	"encoding/json"
	"errors"
	"github.com/go-sanitize/sanitize"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
)

const logMessage = "response:response:"

type HttpStatus int

const STATUS_OK = 200
const STATUS_ERROR = 404

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

func GetEventFromRequest(r *http.Request) (*models.Event, error) {
	eventInput := new(models.ResponseBodyEvent)
	err := json.NewDecoder(r.Body).Decode(eventInput)
	if err != nil {
		return nil, ErrJSONDecoding
	}
	err = ValidateAndSanitize(eventInput)
	if err != nil {
		return nil, err
	}
	result := &models.Event{
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

func GetUserFromRequest(r *http.Request) (*models.User, error) {
	message := logMessage + "GetUserFromRequest:"
	userInput := new(models.ResponseBodyUser)
	jsonReader := strings.NewReader(r.FormValue("json"))
	err := json.NewDecoder(jsonReader).Decode(userInput)
	log.Debug(message + "HERE 2")
	//err := json.NewDecoder(r.Body).Decode(userInput)
	if err != nil {
		return nil, ErrJSONDecoding
	}
	err = ValidateAndSanitize(userInput)
	if err != nil {
		return nil, err
	}
	result := &models.User{
		Name:     userInput.Name,
		Surname:  userInput.Surname,
		Mail:     userInput.Mail,
		Password: userInput.Password,
		About:    userInput.About,
	}
	valid, err := govalidator.ValidateStruct(result)
	log.Debug(message + "HERE 3")
	if err != nil || !valid {
		return nil, ErrValidation
	}
	return result, nil
}

type ResponseBodyEventList struct {
	Events []models.ResponseBodyEvent `json:"events"`
}

func MakeEventForResponse(event *models.Event) models.ResponseBodyEvent {
	return models.ResponseBodyEvent{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Text:        event.Text,
		City:        event.City,
		Category:    event.Category,
		Viewed:      event.Viewed,
		ImgUrl:      event.ImgUrl,
		Tag:         event.Tag,
		Date:        event.Date,
		Geo:         event.Geo,
		AuthorID:    event.AuthorId,
	}
}

func MakeEventListForResponse(events []*models.Event) []models.ResponseBodyEvent {
	result := make([]models.ResponseBodyEvent, len(events))
	for i := 0; i < len(events); i++ {
		result[i] = MakeEventForResponse(events[i])
	}
	return result
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Body    interface{} `json:"body,omitempty"`
}

func ErrorResponse(errorMessage string) *Response {
	return &Response{
		Status:  404,
		Message: errorMessage,
	}
}

func StatusResponse(status HttpStatus) *Response {
	return &Response{
		Status: 200,
		Body:   status,
	}
}

func OkResponse() *Response {
	return &Response{
		Status:  200,
		Message: "OK",
	}
}

func UserResponse(user *models.User) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body: models.ResponseBodyUser{
			ID:      user.ID,
			Name:    user.Name,
			Surname: user.Surname,
			Mail:    user.Mail,
			About:   user.About,
		},
	}
}

type ResponseEventID struct {
	ID string `json:"id"`
}

func EventIdResponse(eventID string) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body: ResponseEventID{
			ID: eventID,
		},
	}
}

func EventsListResponse(events []*models.Event) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body: ResponseBodyEventList{
			Events: MakeEventListForResponse(events),
		},
	}
}

func EventResponse(event *models.Event) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body:    MakeEventForResponse(event),
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

type BaseResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}
