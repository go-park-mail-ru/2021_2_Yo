package response

type statusResponse struct {
	Status string `json:"status"`
	Message string `json:"message,omitempty"`
	Name     string `json:"name,omitempty"`
	Surname  string `json:"surname,omitempty"`
	Mail     string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`	
	Token string `json:"token,omitempty"`	
}


