package main

import (
	"log"

	"github.com/fcproto/prototype/pkg/sensor/dummy"
)

func main() {
	log.Println("starting")

	dummySensor := dummy.NewSensor()

	for i := 0; i < 10; i++ {
		log.Printf("getting value %d: %v", i, dummySensor.GetValues())
	}
}
