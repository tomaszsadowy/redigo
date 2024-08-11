package main

import (
	"fmt"
	"net"
	"strings"
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

	for {
		response := NewResponse(conn)
		value, err := response.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Error... expected array")
		}

		if len(value.array) == 0 {
			fmt.Println("Error... expected array length > 0")
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			// AOF.WRITE(VALUE)
		}

		result := handler(args)
		writer.Write(result)
	}
}
