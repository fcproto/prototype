package main

import (
	"context"
	"fmt"
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
	log.Info("registering sensors...")
	c.RegisterSensor("speed", speed.NewSensor())
	c.RegisterSensor("temperature/env", temperature.New(25, 2), collector.AggregateValues{
		"temperature": collector.AggregatorType_AVG,
	})
	c.RegisterSensor("temperature/track", temperature.New(30, 3), collector.AggregateValues{
		"temperature": collector.AggregatorType_AVG,
	})
	c.RegisterSensor("compass", compass.NewSensor())

	go c.Start()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for {
		select {
		case <-time.After(3 * time.Second):
			fmt.Print(c.Collect())
		case <-ctx.Done():
			log.Info("stopping service...")
			c.Stop()
			return
		}
	}
}
