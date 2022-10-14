package util

type response struct {
	message string
}

func Msg(msg string) any {
	return response{
		message: msg,
	}
}
