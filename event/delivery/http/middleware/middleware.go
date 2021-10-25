package middleware

const logMessage = "event:delivery:http:middleware:"

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}
