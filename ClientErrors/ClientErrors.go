package ClientErrors

import "fmt"

type SendDataError struct{ Err error }
type ReciveDataError struct{ Err error }
type ServerConnectionError struct{ Err error }
type JsonEncodeError struct{ Err error }
type JsonDecodeError struct{ Err error }

type InvalidArgumentCountError struct {
	Arguments uint8
	Expected  uint8
}

func (error *ReciveDataError) Error() string {
	return fmt.Sprintf("error when reciving a response from the server.\n%s", error.Err)
}

func (error *SendDataError) Error() string {
	return fmt.Sprintf("Error when attempting to send the request to the server.\n%s", error.Err)
}

func (error *ServerConnectionError) Error() string {
	return fmt.Sprintf("There has been an error connecting to the server.\nPlease check your connection and try again.\nIf it doesn't work contact the developers and send them this error message:\n\n%s", error.Err)
}

func (error *JsonDecodeError) Error() string {
	return fmt.Sprintf("There has been an Error when attempting to decode the response from the server.\nPlease send this info to the developers:\n%s", error.Err)
}

func (error *JsonEncodeError) Error() string {
	return fmt.Sprintf("There has been an error when attempting to encode the data to be sent to the server.\nPlease send this info to the developers:\n%s", error.Err)
}

func (error *InvalidArgumentCountError) Error() string {
	return fmt.Sprintf("Incorrect number of arguments. got %d, expected %d arguments\nPlease try again", error.Arguments, error.Expected)
}
