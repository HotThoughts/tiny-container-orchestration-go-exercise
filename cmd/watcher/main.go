package main

import (
	"tinygoexercise/pkg/watcher"
)

const port string = "9001"

func main() {
	watcher.InitWatcher(port)
}
