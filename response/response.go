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
	Name        string `json:"name" valid:"type(string)"`
	Description string `json:"description" valid:"type(string)"`
	Viewed      int    `json:"viewed" valid:"type(int)"`
	ImgUrl      string `json:"imgUrl" valid:"type(string)"`
}

type ResponseBodyEventList struct {
	Events []ResponseBodyEvent `json:"events"`
}

func MakeEventListForResponse(events []*models.Event) []ResponseBodyEvent {
	result := make([]ResponseBodyEvent, len(events))
	for i := 0; i < len(events); i++ {
		result[i].Name = events[i].Name
		result[i].Description = events[i].Description
		result[i].Viewed = events[i].Views
		result[i].ImgUrl = events[i].ImgUrl
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

func EventsListResponse(events []*models.Event) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body: ResponseBodyEventList{
			Events: MakeEventListForResponse(events),
		},
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
