package main

import (
	"tinygoexercise/pkg/controller"

	"github.com/jasonlvhit/gocron"
)

const (
	timeInterval int    = 30 // in seconds
	port         string = "9001"
)

var c *controller.Controller

func main() {
	c = controller.New(port)

	s := gocron.NewScheduler()
	//nolint
	s.Every(uint64(timeInterval)).Seconds().Do(task)
	<-s.Start()
}

func task() {
	controller.ControllContainer(c)
}
