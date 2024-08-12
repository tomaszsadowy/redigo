package main

import (
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn, aof *Aof) {
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
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Error... expected array length > 0")
			continue
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
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}

func main() {
	fmt.Println("Listening using port :6379")

	// creating the server
	listen, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listen.Close()

	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	aof.Read(func(val Value) {
		command := strings.ToUpper(val.array[0].bulk)
		args := val.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn, aof)
	}
}
