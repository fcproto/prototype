package main

import (
	"log"

	"github.com/fcproto/prototype/pkg/sensor/compass"
	"github.com/fcproto/prototype/pkg/sensor/dummy"
	"github.com/fcproto/prototype/pkg/sensor/speed"
	"github.com/fcproto/prototype/pkg/sensor/temperature"
)

func main() {
	log.Println("starting")

	dummySensor := dummy.NewSensor()

	for i := 0; i < 10; i++ {
		log.Printf("getting dummy value %d: %v", i, dummySensor.GetValues())
	}

	speedSensor := speed.NewSensor()

	for i := 0; i < 10; i++ {
		log.Printf("getting speed value %d: %v", i, speedSensor.GetValues())
	}

	tempSensor := temperature.NewEnvironmentSensor()

	for i := 0; i < 10; i++ {
		log.Printf("getting temperature value %d: %v", i, tempSensor.GetValues())
	}

	compSensor := compass.NewSensor()

	for i := 0; i < 10; i++ {
		log.Printf("getting compass value %d: %v", i, compSensor.GetValues())
	}
}
