package foopackage

import (
	"github.com/gordonmleigh/golog"
)

var debug = golog.ForPackage(golog.DebugLevel)

// LogSomething just logs something.
func LogSomething() {
	debug.Log("hello from vendor dir")
}
