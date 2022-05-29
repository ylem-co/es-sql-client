package main

type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})

	Debug(...interface{})
	Debugf(string, ...interface{})

	Trace(...interface{})
	Tracef(string, ...interface{})
}
