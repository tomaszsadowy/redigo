package main

import (
	"strconv"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
	"DEL":     del,
	"EXISTS":  exists,
	"INCR":    incr,
	"DECR":    decr,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "Error... wrong no. of arguments for 'set'"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "Error... wrong no. of arguments for 'get'"}
	}

	key := args[0].bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "Error... wrong no. of arguments for 'hset'"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "Error... wrong no. of arguments for 'hget'"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "Error... wrong no. of arguments for 'hgetall'"}
	}

	hash := args[0].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	values := []Value{}
	for k, v := range value {
		values = append(values, Value{typ: "bulk", bulk: k})
		values = append(values, Value{typ: "bulk", bulk: v})
	}

	return Value{typ: "array", array: values}
}

// delete a value
func del(args []Value) Value {
	if len(args) < 1 {
		return Value{typ: "error", str: "Error... wrong no. of arguments for 'del'"}
	}

	deletedCount := 0

	SETsMu.Lock()
	defer SETsMu.Unlock()

	for _, arg := range args {
		key := arg.bulk
		if _, exists := SETs[key]; exists {
			delete(SETs, key)
			deletedCount++
		}

		HSETsMu.Lock()
		if _, exists := HSETs[key]; exists {
			delete(HSETs, key)
			deletedCount++
		}
		HSETsMu.Unlock()
	}

	return Value{typ: "num", num: deletedCount}
}

// check if val exists
func exists(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "Error... wrong no. of arguments for 'exists'"}
	}

	key := args[0].bulk

	SETsMu.RLock()
	defer SETsMu.RUnlock()

	if _, ok := SETs[key]; ok {
		return Value{typ: "num", num: 1}
	}

	return Value{typ: "num", num: 0}
}

func incr(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "Error... wrong no. of arguments for 'incr'"}
	}

	key := args[0].bulk

	SETsMu.Lock()
	defer SETsMu.Unlock()

	value, ok := SETs[key]
	if !ok {
		SETs[key] = "1"
		return Value{typ: "num", num: 1}
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return Value{typ: "error", str: "Error... value is not an integer or is out of range"}
	}

	intValue++
	SETs[key] = strconv.Itoa(intValue)

	return Value{typ: "num", num: intValue}
}

func decr(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "Error... wrong no. of arguments for 'decr'"}
	}

	key := args[0].bulk

	SETsMu.Lock()
	defer SETsMu.Unlock()

	value, ok := SETs[key]
	if !ok {
		SETs[key] = "-1"
		return Value{typ: "num", num: -1}
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return Value{typ: "error", str: "Error... value is not an integer or out of range"}
	}

	intValue--
	SETs[key] = strconv.Itoa(intValue)

	return Value{typ: "num", num: intValue}
}
