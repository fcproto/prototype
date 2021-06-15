package collector

import (
	"sync"

	"github.com/fcproto/prototype/pkg/client/aggregator"
	"github.com/fcproto/prototype/pkg/sensor"
)

// internal struct used for storing the sensor and the required aggregators
type sensorCollector struct {
	sensorMutex sync.Mutex
	sensor      sensor.Sensor
	aggregators map[string]*aggregator.Aggregator
}

func newSensorCollector(sensor sensor.Sensor, aggValues []AggregateValues, aggregatorSize int) *sensorCollector {
	aggregators := make(map[string]*aggregator.Aggregator)
	for _, aggVal := range aggValues {
		for v, t := range aggVal {
			aggregators[v] = aggregator.New(t, aggregatorSize)
		}
	}
	return &sensorCollector{
		sensor:      sensor,
		aggregators: aggregators,
	}
}

func (cs *sensorCollector) aggregate() {
	cs.sensorMutex.Lock()
	defer cs.sensorMutex.Unlock()
	if len(cs.aggregators) == 0 {
		return
	}
	values := cs.sensor.GetValues()
	for k, val := range values {
		if _, ok := cs.aggregators[k]; !ok {
			// we just want to keep the last value, if no aggregator was defined
			cs.aggregators[k] = aggregator.NewSingle()
		}
		cs.aggregators[k].AddValue(val)
	}
	for k, agg := range cs.aggregators {
		if val, ok := values[k]; ok {
			agg.AddValue(val)
		}
	}
}

func (cs *sensorCollector) getValues() sensor.Values {
	cs.sensorMutex.Lock()
	defer cs.sensorMutex.Unlock()
	if len(cs.aggregators) == 0 {
		//if no aggregators are enabled, get the current values
		return cs.sensor.GetValues()
	}
	values := make(sensor.Values)
	// fetch the values from the aggregators
	for k, agg := range cs.aggregators {
		values[k] = agg.Aggregate()
	}
	return values
}
