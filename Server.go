package main

import (
	"CloudDrive/server/RequestHandlers"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), "is connected")
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
}

func main() {
	login_manager, err := RequestHandlers.InitializeFactory()
	if err != nil {
		fmt.Println(err.Error())
	}
	errs := login_manager.Signup("sdfds", "ssdfdsjdsfdfdfs", "dfdf@jdfj.com")
	if len(errs) > 0 {
		fmt.Println(errs)
	} else {
		fmt.Println(login_manager.GetLoggedUsers())
	}
	dd, err := RequestHandlers.GetManager()
	if err == nil {
		fmt.Println(dd.GetLoggedUsers())
	}

	addr := "192.168.50.191:12345"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server is listening on %s...\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}

}
