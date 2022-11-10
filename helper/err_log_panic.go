package helper

import "log"

func ErrLogPanic(err error) {
	if err != nil {
		log.Panic(err)
	}
}
