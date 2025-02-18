package core

import (
	"errors"
	"io"
	"strconv"
	"time"
)

var NIL_RESP = []byte("$-1\r\n")

func evalPING(args []string, c io.ReadWriter) error {
	if len(args) >= 2 {
		return errors.New("ERR wrong number of arguments for 'ping' command")
	}

	var b []byte
	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}
	_, err := c.Write(b)
	return err
}

func evalSET(args []string, c io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("ERR wrong number of arguments for 'set' command")
	}

	key := args[0]
	value := args[1]

	ttlMillis := int64(-1)
	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			{
				// check if there is a value after EX
				if i+1 == len(args) {
					return errors.New("ERR syntax error")
				}
				i++

				// check if the value is a number
				ttlSeconds, err := strconv.ParseInt(args[i], 10, 64)
				if err != nil {
					return errors.New("ERR value is not an integer or out of range")
				}

				ttlMillis = ttlSeconds * 1000
			}
		}
	}

	obj := NewObj(value, ttlMillis)
	Set(key, obj)
	_, err := c.Write([]byte("+OK\r\n"))
	if err != nil {
		return errors.New("ERR failed to write response")
	}
	return nil
}

func evalGET(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("ERR wrong number of arguments for 'get' command")
	}

	key := args[0]
	obj := Get(key)
	if obj == nil {
		_, err := c.Write(NIL_RESP)
		if err != nil {
			return errors.New("ERR failed to write response")
		}
		return nil
	}

	if obj.expireAt > 0 && obj.expireAt < time.Now().UnixMilli() {
		_, err := c.Write(NIL_RESP)
		if err != nil {
			return errors.New("ERR failed to write response")
		}
		return nil
	}

	b := Encode(obj.value, false)
	_, err := c.Write(b)
	if err != nil {
		return errors.New("ERR failed to write response")
	}
	return nil
}

func evalTTL(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("ERR wrong number of arguments for 'ttl' command")
	}

	key := args[0]
	obj := Get(key)
	if obj == nil {
		_, err := c.Write([]byte(":-2\r\n"))
		if err != nil {
			return errors.New("ERR failed to write response")
		}
		return nil
	}

	// if key does not have an expiry time
	if obj.expireAt == -1 {
		_, err := c.Write([]byte(":-1\r\n"))
		if err != nil {
			return errors.New("ERR failed to write response")
		}
		return nil
	}

	// compute the time remaining for the key to expire
	durationMS := obj.expireAt - time.Now().UnixMilli()

	// if the key has already expired
	if durationMS < 0 {
		_, err := c.Write([]byte(":-2\r\n"))
		if err != nil {
			return errors.New("ERR failed to write response")
		}
		return nil
	}

	b := Encode(durationMS/1000, false)
	_, err := c.Write(b)
	if err != nil {
		return errors.New("ERR failed to write response")
	}
	return nil
}

func EvalAndRespond(cmd *RedisCmd, c io.ReadWriter) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd.Args, c)
	case "SET":
		return evalSET(cmd.Args, c)
	case "GET":
		return evalGET(cmd.Args, c)
	case "TTL":
		return evalTTL(cmd.Args, c)
	default:
		return evalPING(cmd.Args, c)
	}
}
