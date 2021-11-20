package response

import (
	models2 "backend/pkg/models"
)

const logMessage = "response:response:"

type HttpStatus int

const STATUS_OK = 200
const STATUS_ERROR = 404

type BaseResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Body    interface{} `json:"body,omitempty"`
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

func ErrorResponse(errorMessage string) *Response {
	return &Response{
		Status:  404,
		Message: errorMessage,
	}
}

func UserResponse(user *models2.User) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body:    MakeUserResponseBody(user),
	}
}

func UserListResponse(users []*models2.User) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body:    MakeUserListResponseBody(users),
	}
}

func EventResponse(event *models2.Event) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body:    MakeEventResponseBody(event),
	}
}

func EventIdResponse(eventID string) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body: models2.EventIDResponseBody{
			ID: eventID,
		},
	}
}

func EventListResponse(events []*models2.Event) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body:    MakeEventListResponseBody(events),
	}
}
