package main

import (
	"github.com/fcproto/prototype/pkg/logger"
	"github.com/fcproto/prototype/pkg/sensor/compass"
	"github.com/fcproto/prototype/pkg/sensor/dummy"
	"github.com/fcproto/prototype/pkg/sensor/speed"
	"github.com/fcproto/prototype/pkg/sensor/temperature"
)

func main() {
	log := logger.New("sensor-test")
	log.Info("starting")

	dummySensor := dummy.NewSensor()

	for i := 0; i < 10; i++ {
		log.Infof("getting dummy value %d: %v", i, dummySensor.GetValues())
	}

	speedSensor := speed.NewSensor()

	for i := 0; i < 10; i++ {
		log.Infof("getting speed value %d: %v", i, speedSensor.GetValues())
	}

	tempSensor := temperature.NewEnvironmentSensor()

	for i := 0; i < 10; i++ {
		log.Infof("getting temperature value %d: %v", i, tempSensor.GetValues())
	}

	compSensor := compass.NewSensor()

	for i := 0; i < 10; i++ {
		log.Infof("getting compass value %d: %v", i, compSensor.GetValues())
	}
}
