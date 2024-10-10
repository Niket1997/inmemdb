package core

import "errors"

// reads RESP encoded simple string from data & returns
// the string, the delta & the optional error
func readSimpleString(data []byte) (string, int, error) {
	// first character of simple string is '+'
	pos := 1

	for ; data[pos] != '\r'; pos++ {
	}

	return string(data[1:pos]), pos + 2, nil
}

// reads RESP encoded error from data & returns
// the error string, the delta & the optional error
func readError(data []byte) (string, int, error) {
	// first character for error is '-'
	// first character for simple string is '+'
	// that's the only difference RESP diff
	// between simple string & error
	return readSimpleString(data)
}

// reads RESP encoded int64 from data & returns
// the int64, the delta & the optional error
func readInt64(data []byte) (int64, int, error) {
	// first character of an integer is ':'
	pos := 1
	var value int64 = 0

	for ; data[pos] != '\r'; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}

	return value, pos + 2, nil
}

// reads RESP encoded bulk string & returns
// the string, the delta & optional error
func readBulkString(data []byte) (string, int, error) {
	// first character of bulk string is '$'
	pos := 1

	// read the length & forward the pos by
	// the length of integer + the first special character + 2 (CRLF)
	length, delta := readLength(data[pos:])
	pos += delta

	return string(data[pos : pos+length]), pos + length + 2, nil
}

// reads the length typically the first integer of the string
// until hit by a non-digit byte & return the integer length
// & the delta = length + 2 (CRLF, i.e. \r\n)
func readLength(data []byte) (int, int) {
	pos := 0
	length := 0

	for pos = range data {
		b := data[pos]
		if !(b >= '0' && b <= '9') {
			return length, pos + 2
		}
		length = length*10 + int(b-'0')
	}

	return 0, 0
}

// reads RESP encoded array from data & returns
// the array, the delta & optional error
func readArray(data []byte) (interface{}, int, error) {
	// first character of array is '*'
	pos := 1

	// read the length & forward the pos by
	// the length of integer + the first special character + 2 (CRLF)
	count, delta := readLength(data[pos:])
	pos += delta

	var elems []interface{} = make([]interface{}, count)

	for i := range elems {
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}

		elems[i] = elem
		pos += delta
	}

	return elems, pos, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt64(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}

	return nil, 0, nil
}

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}

	value, _, err := DecodeOne(data)
	return value, err
}
