package response

import (
	"backend/pkg/models"
)

const logMessage = "response:response:"

type HttpStatus int

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

func UserResponse(user *models.User) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body:    MakeUserResponseBody(user),
	}
}

func UserListResponse(users []*models.User) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body:    MakeUserListResponseBody(users),
	}
}

func EventResponse(event *models.Event) *Response {
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
		Body: models.EventIDResponseBody{
			ID: eventID,
		},
	}
}

func EventListResponse(events []*models.Event) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body:    MakeEventListResponseBody(events),
	}
}

func SubscribedResponse(result bool) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body: models.SubscribedResponseBody{
			Result: result,
		},
	}
}

func FavouriteResponse(result bool) *Response {
	return &Response{
		Status:  200,
		Message: "",
		Body: models.FavouriteResponseBody{
			Result: result,
		},
	}
}
