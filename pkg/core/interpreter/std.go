package main

import "time"

type timestampFunc struct{}

func (f timestampFunc) call(sc *Scope, args ...interface{}) (interface{}, error) {
	return time.Now().Unix(), nil
}
