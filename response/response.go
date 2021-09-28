package response

type ResponseBody struct {
	Message  string `json:"message,omitempty"`
	Name     string `json:"name,omitempty"`
	Surname  string `json:"surname,omitempty"`
	Mail     string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

type Response struct {
	Status   int    `json:"status"`
	Body ResponseBody `json:"body"`
}


func ErrorResponse(errorMessage string) *Response {
	return &Response{
		Status:  404,
		Body :ResponseBody {
			Message: errorMessage,
		},
	}
}

func OkResponse() *Response {
	return &Response{
		Status: 200,
		Body :ResponseBody {
			
		},
	}
}

func UsernameResponse(name string) *Response {
	return &Response{
		Status: 200,
		Body :ResponseBody {
			Name :name,
		},
	}
}
