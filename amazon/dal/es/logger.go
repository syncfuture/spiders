package es

import (
	log "github.com/kataras/golog"
)

type logger struct {
	Type string
}

func (x *logger) Printf(format string, v ...interface{}) {
	switch x.Type {
	case "debug":
		log.Debugf(format, v...)
		break
	case "info":
		log.Infof(format, v...)
		break
	case "error":
		log.Errorf(format, v...)
		break
	}
}
