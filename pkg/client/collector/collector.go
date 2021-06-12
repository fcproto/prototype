package collector

import (
	"sync"
	"time"

	"github.com/fcproto/prototype/pkg/client/api"
	"github.com/fcproto/prototype/pkg/logger"
	"github.com/fcproto/prototype/pkg/sensor"
	"github.com/ipfs/go-log/v2"
)

type AggregateValues map[string]AggregatorType

// internal struct used for storing the sensor and the required aggregators
type collectSensor struct {
	sensor      sensor.Sensor
	aggregators map[string]*Aggregator
}

func (cs *collectSensor) refresh() {
	if len(cs.aggregators) == 0 {
		return
	}
	values := cs.sensor.GetValues()
	for k, agg := range cs.aggregators {
		if val, ok := values[k]; ok {
			agg.AddValue(val)
		}
	}
}

func (cs *collectSensor) getValues() sensor.Values {
	values := cs.sensor.GetValues()
	//replace current values with aggregated values
	for k, agg := range cs.aggregators {
		if _, ok := values[k]; !ok {
			// skip if aggregator is not in values
			continue
		}
		// add the current value to the aggregator and aggregate
		agg.AddValue(values[k])
		values[k] = agg.Aggregate()
	}
	return values
}

type Collector struct {
	log        *log.ZapEventLogger
	interval   time.Duration
	tickerDone chan struct{}

	sensorsMutex sync.Mutex
	sensors      map[string]*collectSensor
}

func New() *Collector {
	return &Collector{
		log:        logger.New("collector"),
		interval:   time.Second,
		tickerDone: make(chan struct{}),
		sensors:    make(map[string]*collectSensor),
	}
}

func (c *Collector) RegisterSensor(name string, sensor sensor.Sensor, agg ...AggregateValues) {
	c.sensorsMutex.Lock()
	defer c.sensorsMutex.Unlock()

	aggregators := make(map[string]*Aggregator)
	for _, aggEl := range agg {
		for v, t := range aggEl {
			aggregators[v] = NewAggregator(t, 5)
		}
	}
	c.log.Debugf("adding new sensor: %s", name)
	c.sensors[name] = &collectSensor{
		sensor:      sensor,
		aggregators: aggregators,
	}
}

func (c *Collector) Collect() *api.SensorData {
	c.sensorsMutex.Lock()
	defer c.sensorsMutex.Unlock()
	res := &api.SensorData{
		Timestamp: time.Now(),
		Sensors:   make(map[string]sensor.Values),
	}
	for k, v := range c.sensors {
		res.Sensors[k] = v.getValues()
	}
	return res
}

func (c *Collector) refreshSensors() {
	c.sensorsMutex.Lock()
	defer c.sensorsMutex.Unlock()
	//c.log.Debug("refreshing sensors...")
	for _, v := range c.sensors {
		v.refresh()
	}
}

func (c *Collector) Start() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.refreshSensors()
		case <-c.tickerDone:
			c.log.Debug("refreshing stopped")
			return
		}
	}
}

func (c *Collector) Stop() {
	c.tickerDone <- struct{}{}
}
