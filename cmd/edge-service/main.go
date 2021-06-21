package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fcproto/prototype/pkg/client"
	"github.com/fcproto/prototype/pkg/client/aggregator"
	"github.com/fcproto/prototype/pkg/client/collector"
	"github.com/fcproto/prototype/pkg/logger"
	"github.com/fcproto/prototype/pkg/sensor/compass"
	"github.com/fcproto/prototype/pkg/sensor/gps"
	"github.com/fcproto/prototype/pkg/sensor/temperature"
)

func main() {
	log := logger.New("main")
	log.Info("starting edge service...")
	service, err := client.NewService("https://envwzzmqa85j.x.pipedream.net/", 120)
	if err != nil {
		log.Fatal(err)
	}

	c := collector.New()
	log.Info("registering sensors...")
	c.RegisterSensor("gps", gps.NewSensor())
	c.RegisterSensor("compass", compass.NewSensor())
	c.RegisterSensor("temperature/env", temperature.New(25, 2), collector.AggregateValues{
		"temperature": aggregator.TypeAvg,
	})
	c.RegisterSensor("temperature/track", temperature.New(30, 3), collector.AggregateValues{
		"temperature": aggregator.TypeAvg,
	})

	log.Info("starting collector...")
	go c.Start()
	defer c.Stop()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		for range time.Tick(time.Second) {
			log.Debug("collecting and saving data...")
			if err := service.SubmitSensorData(c.Collect()); err != nil {
				log.Error(err)
			}
		}
	}()

	go func() {
		for range time.Tick(10 * time.Second) {
			log.Debug("syncing...")
			if err := service.Sync(); err != nil {
				log.Error(err)
			}
		}
	}()

	<-ctx.Done()
	log.Info("stopping service...")
}
