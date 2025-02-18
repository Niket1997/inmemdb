package core

import "time"

var store map[string]*Obj

type Obj struct {
	value    interface{}
	expireAt int64
}

func init() {
	store = make(map[string]*Obj)
}

func NewObj(value interface{}, ttlMillis int64) *Obj {
	expireAt := int64(-1)
	if ttlMillis > 0 {
		expireAt = time.Now().UnixMilli() + ttlMillis
	}
	return &Obj{value: value, expireAt: expireAt}
}

func Set(key string, value *Obj) {
	store[key] = value
}

func Get(key string) *Obj {
	return store[key]
}
