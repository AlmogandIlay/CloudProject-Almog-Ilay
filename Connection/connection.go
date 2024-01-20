package connection

import (
	"fmt"
	"net"
)

var Socket *net.Conn

const (
	ip_addr = "192.168.50.191:12345"
)

func Connect() {
	var err error
	*Socket, err = net.Dial("tcp", ip_addr)
	if err != nil {
		fmt.Println("There has been an error connecting to the server.\nPlease check your connection and try again.\nIf it doesn't work contact the developers and send them this error message:\n", err.Error())
	}
}
