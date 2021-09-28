package response

type Response struct {
	Status   int    `json:"status"`
	Message  string `json:"message,omitempty"`
	Name     string `json:"name,omitempty"`
	Surname  string `json:"surname,omitempty"`
	Mail     string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

func ErrorResponse(errorMessage string) *Response {
	return &Response{
		Status:  500,
		Message: errorMessage,
	}
}

func OkResponse() *Response {
	return &Response{
		Status: 200,
	}
}

func UsernameResponse(name string) *Response {
	return &Response{
		Status: 200,
		Name:   name,
	}
}
