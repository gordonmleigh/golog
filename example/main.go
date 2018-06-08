package main

import "github.com/gordonmleigh/golog"

func main() {
	debug := golog.NewLogger("DEBUG:hello")
	golog.SetWriter(golog.MustParseNameFilter("DEBUG:*"), golog.ConsoleWriter)
	golog.SetWriter(golog.MustParseNameFilter("ERROR:*"), golog.ConsoleWriter)
	err := golog.NewLogger("ERROR:hello")
	debug.Log("hello world")
	err.Log(
		"BANG!",
		golog.Val("value1", 5),
		golog.Val("translation", "cinq"),
	)
}
