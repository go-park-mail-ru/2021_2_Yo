package response

import (
	log "backend/logger"
	"backend/models"
	"encoding/json"
	"net/http"
)

const logMessage = "response:response:"

const STATUS_OK = 200
const STATUS_ERROR = 404

type ResponseBodyUser struct {
	Name     string `json:"name,omitempty" valid:"type(string)"`
	Surname  string `json:"surname,omitempty" valid:"type(string)"`
	Mail     string `json:"email,omitempty" valid:"email"`
	Password string `json:"password,omitempty" valid:"type(string)"`
}

type ResponseBodyEvent struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Text        string   `json:"text"`
	City        string   `json:"city"`
	Category    string   `json:"category"`
	Viewed      int      `json:"viewed"`
	Tag         []string `json:"tag"`
	Date        string   `json:"date"`
	Geo         string   `json:"geo"`
}

type ResponseBodyEventList struct {
	Events []ResponseBodyEvent `json:"events"`
}

func MakeEventForResponse(event *models.Event) ResponseBodyEvent {
	return ResponseBodyEvent{
		Title:       event.Title,
		Description: event.Description,
		Text:        event.Text,
		City:        event.City,
		Category:    event.Category,
		Viewed:      event.Viewed,
		Tag:         event.Tag,
		Date:        event.Date,
		Geo:         event.Geo,
	}
}

func MakeEventListForResponse(events []*models.Event) []ResponseBodyEvent {
	result := make([]ResponseBodyEvent, len(events))
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

func OkResponse() *Response {
	return &Response{
		Status:  200,
		Message: "OK",
	}
}

func UsernameResponse(name string) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body: ResponseBodyUser{
			Name: name,
		},
	}
}

type ResponseEventID struct {
	ID string `json:"id"`
}

func EventIdResponse(eventID string) *Response {
	return &Response{
		Status:  0,
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
	log.Debug(message+"response to send =", string(b))
	w.Write(b)
}

//For docs
type BaseResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}
