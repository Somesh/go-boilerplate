package safe

import (
	"errors"
	"log"
	"runtime/debug"

	"github.com/Somesh/go-boilerplate/tools/panics"
)

func Recover() {
	var err error
	r := recover()

	if r != nil {
		switch v := r.(type) {
		case string:
			err = errors.New(v)
		case error:
			err = v
		default:
			err = errors.New("unknown error")
		}
		log.Printf("Recovered Panic, Error: %+v, Stack:%+v", string(debug.Stack()))

		panics.Capture(err.Error(), string(debug.Stack()))
	}
}
