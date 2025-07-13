package main

var verbose logger

type logger struct {
	logf func(format string, v ...any)
}

func (l logger) Printf(format string, v ...any) {
	if l.logf == nil {
		return
	}
	l.logf(format, v...)
}
