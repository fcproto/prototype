package collector

import (
	"sync"
)

type AggregatorType byte

const (
	AggregatorType_AVG AggregatorType = iota
	AggregatorType_MIN AggregatorType = iota
	AggregatorType_MAX AggregatorType = iota
)

type Aggregator struct {
	aggType     AggregatorType
	index       int
	valuesMutex sync.RWMutex
	values      []float64
	filled      bool
}

func NewAggregator(aggType AggregatorType, size int) *Aggregator {
	return &Aggregator{
		aggType: aggType,
		index:   0,
		values:  make([]float64, size),
		filled:  false,
	}
}

func (a *Aggregator) AddValue(v float64) {
	a.valuesMutex.Lock()
	defer a.valuesMutex.Unlock()
	a.values[a.index] = v
	a.index++
	if a.index == len(a.values) {
		a.index = 0
		// after filling the buffer for the first time, we don't need to check how many values should be aggregated
		a.filled = true
	}
}

// arithmetic mean
func (a *Aggregator) aggregateAvg() float64 {
	a.valuesMutex.RLock()
	defer a.valuesMutex.RUnlock()

	// determine size
	size := len(a.values)
	if !a.filled {
		// we don't have any values, so just return null
		if a.index == 0 {
			return 0
		}
		size = a.index
	}

	agg := a.values[0]
	for i := 1; i < size; i++ {
		agg += a.values[i]
	}

	return agg / float64(size)
}

func (a *Aggregator) Aggregate() float64 {
	switch a.aggType {
	case AggregatorType_MIN:
		// TODO
		return 0
	case AggregatorType_MAX:
		// TODO
		return 0
	}

	// default is avg
	return a.aggregateAvg()
}
