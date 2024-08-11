
<h1 align="center">
  <img src="https://github.com/user-attachments/assets/f093a8bc-2be5-4901-b0e7-6d087016e438" width=250>
</h1>

<p align="center">
  <i>A lightweight in-memory key-value store, written in Go 🌀</i>
</p> 

## Introduction

**Redigo** (redis - go) is a simple yet powerful in-memory key-value store written in *Go*. It is designed to provide a subset of *Redis*-like functionality with a focus on simplicity and ease of use. Redigo supports a variety of commands to interact with the stored data and ensures data persistence using an append-only file (AOF). This project is designed for educational purposes and demonstrates how to build a TCP server that can handle Redis-like commands.

## Details
The server listens on port 6379 and accepts incoming TCP connections. It reads commands from the client, processes them, and sends back the response.
Command Handlers: Each supported command has a corresponding handler function defined in handler.go. These handlers are registered in a map for easy lookup.
The server supports basic AOF persistence. Commands that modify the data store (SET, HSET) are appended to an AOF file (database.aof). On startup, the server replays the AOF file to restore the data store.

## Installation
1. Clone the git repository using:
   ```sh
   git clone https://github.com/tomaszsadowy/redigo.git
   ```
2. Install redis and the redis-stack (this is for our client):
   
   https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/
   
4. Test test the client by using the following commands:
   ```sh
   $ redis-cli
   redis> PING

   Output: PONG
   ```

## Usage

1. Start the Server:
   ```sh
   $ go run *.go
   ```
2. Connect using redis cli:
   ```sh
   $ redis-cli
   ```
3. Execute commands:
   ```sh
    redis> PING
    PONG
   
    redis> SET mykey hello
    OK
   
    redis> GET mykey
    "hello"

    redis DEL mykey
    1

    redis EXISTS mykey
    0
   
    redis> HSET users u1 tomasz
    1
   
    redis> HGET users u1
    "tomasz"
   
    redis> HGETALL users
    1) "u1"
    2) "tomasz"
   ```
4. Notes:
   There are more commands, such as `INCR`, `DECR`, and more will be added continuously :)
