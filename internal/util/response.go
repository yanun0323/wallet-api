package util

type Response struct {
	Message string `json:"message"`
}

func Msg(msg string) Response {
	return Response{
		Message: msg,
	}
}
