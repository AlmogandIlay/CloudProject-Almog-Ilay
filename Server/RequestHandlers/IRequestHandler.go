package requesthandlers

type IRequestHandler interface {
	isValidRequest() bool
	handleRequest()
	error()
}
