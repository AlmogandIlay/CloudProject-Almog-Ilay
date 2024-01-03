package RequestHandlers

type ResponeType int

const (
	ErrorRespone  ResponeType = iota // 0
	LoginRespone  ResponeType = iota // 1
	SignupRespone ResponeType = iota // 2
)

type ResponeInfo struct {
	Respone    string
	NewHandler *IRequestHandler
}

func buildRespone(respone string, handler *IRequestHandler) ResponeInfo {
	return ResponeInfo{Respone: respone, NewHandler: handler}
}
