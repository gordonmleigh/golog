package main

import "github.com/gordonmleigh/golog"

func main() {
	debug := golog.NewLogger(golog.DebugLevel, "hello")

	golog.SetWriter(golog.ConsoleWriter, golog.DebugLevel, golog.Wildcard)
	golog.SetWriter(golog.ConsoleWriter, "ERROR:*")

	err := golog.NewLogger("ERROR:hello")

	debug.Log("hello world")

	err.Log(
		"BANG!",
		golog.Val("value1", 5),
		golog.Val("translation", "cinq"),
	)
}
