package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Listening using port :6379")

	// creating the server
	listen, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// listening to connections
	conn, err := listen.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()
}
