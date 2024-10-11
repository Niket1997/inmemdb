package core_test

import (
	"fmt"
	"testing"

	"github.com/Niket1997/inmemdb/core"
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

func TestDecodeArrayString(t *testing.T) {
	cases := map[string][]string{
		"*1\r\n+PING\r\n":                     {"PING"},
		"*1\r\n$4\r\nPING\r\n":                {"PING"},
		"*2\r\n$4\r\nPING\r\n$5\r\nhello\r\n": {"PING", "hello"},
	}

	for k, v := range cases {
		values, err := core.DecodeArrayString([]byte(k))
		if err != nil {
			t.Fail()
		}

		if len(values) != len(v) {
			t.Fail()
		}

		for i := range values {
			if values[i] != v[i] {
				t.Fail()
			}
		}
	}
}

func TestEncode(t *testing.T) {
	type testCase struct {
		value    interface{}
		isSimple bool
		expected []byte
	}

	cases := []testCase{
		{
			"PONG",
			true,
			[]byte(fmt.Sprintf("+PONG\r\n")),
		},
		{
			"PONG",
			false,
			[]byte(fmt.Sprintf("$4\r\nPONG\r\n")),
		},
	}

	for i := range cases {
		test := cases[i]
		v := core.Encode(test.value, test.isSimple)
		if len(v) != len(test.expected) {
			t.Fail()
		}

		for j := range v {
			if v[j] != test.expected[j] {
				t.Fail()
			}
		}
	}
}
