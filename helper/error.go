package helper

import (
	log "github.com/sirupsen/logrus"
)

// FatalIfError throws fatal if error
func FatalIfError(err error, extra ...interface{}) {
	if err != nil {
		a := []interface{}{err}
		for i := range extra {
			a = append(a, ": ", extra[i])
		}
		log.Fatal(a...)
	}
}

// ErrorIfError shows error if error
func ErrorIfError(err error, extra ...interface{}) {
	if err != nil {
		a := []interface{}{err}
		for i := range extra {
			a = append(a, ": ", extra[i])
		}
		log.Error(a...)
	}
}
