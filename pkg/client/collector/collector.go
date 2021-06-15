package collector

import (
	"sync"
	"time"

	"github.com/fcproto/prototype/pkg/api"
	"github.com/fcproto/prototype/pkg/client/aggregator"
	"github.com/fcproto/prototype/pkg/logger"
	"github.com/fcproto/prototype/pkg/sensor"
	"github.com/ipfs/go-log/v2"
)

type AggregateValues map[string]aggregator.Type

type Collector struct {
	log        *log.ZapEventLogger
	interval   time.Duration
	tickerDone chan struct{}

	sensorsMutex   sync.Mutex
	sensors        map[string]*sensorCollector
	aggregatorSize int
}

func New() *Collector {
	return &Collector{
		log:            logger.New("collector"),
		interval:       time.Millisecond * 50,
		tickerDone:     make(chan struct{}),
		sensors:        make(map[string]*sensorCollector),
		aggregatorSize: 100,
	}
}

func (c *Collector) RegisterSensor(name string, sensor sensor.Sensor, aggValues ...AggregateValues) {
	c.sensorsMutex.Lock()
	defer c.sensorsMutex.Unlock()
	c.log.Debugf("adding new sensor: %s", name)
	c.sensors[name] = newSensorCollector(sensor, aggValues, c.aggregatorSize)
}

func (c *Collector) Collect() *api.SensorData {
	c.sensorsMutex.Lock()
	defer c.sensorsMutex.Unlock()
	res := api.NewSensorData()
	for k, v := range c.sensors {
		res.Sensors[k] = v.getValues()
	}
	return res
}

func (c *Collector) aggregateSensors() {
	c.sensorsMutex.Lock()
	defer c.sensorsMutex.Unlock()
	//c.log.Debug("aggregating sensors...")
	for _, v := range c.sensors {
		v.aggregate()
	}
}

func (c *Collector) Start() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.aggregateSensors()
		case <-c.tickerDone:
			c.log.Debug("aggregating stopped")
			return
		}
	}
}

func (c *Collector) Stop() {
	c.tickerDone <- struct{}{}
}
