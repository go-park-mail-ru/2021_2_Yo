package response

import (
	models2 "backend/internal/models"
)

const logMessage = "response:response:"

type HttpStatus int

type Response struct {
	Status  HttpStatus  `json:"status"`
	Message string      `json:"message,omitempty"`
	Body    interface{} `json:"body,omitempty"`
}

func StatusResponse(status HttpStatus) *Response {
	return &Response{
		Status: status,
	}
}

func OkResponse() *Response {
	return &Response{
		Status:  200,
		Message: "OK",
	}
}

func UserResponse(user *models2.User) *Response {
	return &Response{
		Status: 200,
		Body:   MakeUserResponseBody(user),
	}
}

func UserListResponse(users []*models2.User) *Response {
	return &Response{
		Status: 200,
		Body:   MakeUserListResponseBody(users),
	}
}

func EventResponse(event *models2.Event) *Response {
	return &Response{
		Status: 200,
		Body:   MakeEventResponseBody(event),
	}
}

func EventIdResponse(eventID string) *Response {
	return &Response{
		Status: 200,
		Body: models2.EventIDResponseBody{
			ID: eventID,
		},
	}
}

func EventListResponse(events []*models2.Event) *Response {
	return &Response{
		Status: 200,
		Body:   MakeEventListResponseBody(events),
	}
}

func SubscribedResponse(result bool) *Response {
	return &Response{
		Status: 200,
		Body: models2.SubscribedResponseBody{
			Result: result,
		},
	}
}

func FavouriteResponse(result bool) *Response {
	return &Response{
		Status: 200,
		Body: models2.FavouriteResponseBody{
			Result: result,
		},
	}
}

func CitiesResponse(cities []string) *Response {
	return &Response{
		Status: 200,
		Body: models2.CitiesResponseBody{
			Cities: cities,
		},
	}
}
