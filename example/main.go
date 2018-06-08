package main

import (
	"foopackage"

	"github.com/gordonmleigh/golog"
)

var debug = golog.ForPackage(golog.DebugLevel)

func main() {
	golog.SetWriter(golog.ConsoleWriter, golog.DebugLevel, golog.Wildcard)
	golog.SetWriter(golog.ConsoleWriter, "ERROR:*")

	err := golog.NewLogger(golog.ErrorLevel, "hello")

	debug.Log("hello world")
	debug.Log(golog.GetPackageName(-1))
	foopackage.LogSomething()

	err.Log(
		"BANG!",
		golog.Val("value1", 5),
		golog.Val("translation", "cinq"),
	)
}
