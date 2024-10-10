package core_test

import (
	"fmt"
	"github.com/Niket1997/inmemdb/core"
	"testing"
)

func TestSimpleStringDecode(t *testing.T) {
	cases := map[string]string{
		"+OK\r\n":          "OK",
		"+Hello World\r\n": "Hello World",
	}

	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		if value != v {
			t.Fail()
		}
	}
}

func TestErrorDecode(t *testing.T) {
	cases := map[string]string{
		"-Error Message\r\n": "Error Message",
	}

	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		if value != v {
			t.Fail()
		}
	}
}

func TestInt64Decode(t *testing.T) {
	cases := map[string]int64{
		":0\r\n":    0,
		":1100\r\n": 1100,
	}

	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		if value != v {
			t.Fail()
		}
	}
}

func TestBulkStringDecode(t *testing.T) {
	cases := map[string]string{
		"$0\r\n\r\n":      "",
		"$5\r\nhello\r\n": "hello",
	}

	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		if value != v {
			t.Fail()
		}
	}
}

func TestArrayDecode(t *testing.T) {
	cases := map[string][]interface{}{
		"*0\r\n":                               {},
		"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n": {"hello", "world"},
		"*3\r\n:1\r\n:2\r\n:3\r\n":             {int64(1), int64(2), int64(3)},
		"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$11\r\nHello World\r\n":     {int64(1), int64(2), int64(3), int64(4), "Hello World"},
		"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-Error\r\n": {[]int64{int64(1), int64(2), int64(3)}, []string{"Hello", "Error"}},
	}

	for k, v := range cases {
		value, _ := core.Decode([]byte(k))

		array := value.([]interface{})

		if len(array) != len(v) {
			t.Fail()
		}

		for i := range array {
			if fmt.Sprintf("%v", v[i]) != fmt.Sprintf("%v", array[i]) {
				t.Fail()
			}
		}
	}
}
