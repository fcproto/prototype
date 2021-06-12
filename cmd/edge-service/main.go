package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fcproto/prototype/pkg/client/collector"
	"github.com/fcproto/prototype/pkg/logger"
	"github.com/fcproto/prototype/pkg/sensor/compass"
	"github.com/fcproto/prototype/pkg/sensor/speed"
	"github.com/fcproto/prototype/pkg/sensor/temperature"
)

func main() {
	log := logger.New("main")
	log.Info("starting edge service...")

	c := collector.New()
	c.RegisterSensor("speed", speed.NewSensor())
	c.RegisterSensor("temperature/env", temperature.NewEnvironmentSensor(), collector.AggregateValues{
		"environment": collector.AggregatorType_AVG,
	})
	c.RegisterSensor("temperature/track", temperature.NewTrackSensor(), collector.AggregateValues{
		"track": collector.AggregatorType_AVG,
	})

	c.RegisterSensor("compass", compass.NewSensor())

	go c.Start()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for {
		select {
		case <-time.After(7 * time.Second):
			log.Debug(c.Collect())
		case <-ctx.Done():
			log.Info("stopping service...")
			return
		}
	}
}
