package collector

import (
	"fmt"
	"sync"
)

type AggregatorType byte

const (
	AggregatorType_AVG AggregatorType = iota
	AggregatorType_MIN AggregatorType = iota
	AggregatorType_MAX AggregatorType = iota
)

func (at AggregatorType) String() string {
	switch at {
	case AggregatorType_AVG:
		return "AVG"
	case AggregatorType_MIN:
		return "MIN"
	case AggregatorType_MAX:
		return "MAX"
	default:
		return "UNKNOWN"
	}
}

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

func (a *Aggregator) String() string {
	a.valuesMutex.Lock()
	defer a.valuesMutex.Unlock()
	return fmt.Sprintf("{type=%s, index=%d, filled=%t, values=%v}", a.aggType.String(), a.index, a.filled, a.values)
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

func reducerSum(a, b float64) float64 {
	return a + b
}

func reducerMin(a, b float64) float64 {
	if a > b {
		return b
	} else {
		return a
	}
}

func reducerMax(a, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func (a *Aggregator) aggregateFn(reducer func(float64, float64) float64, divBySize bool) float64 {
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
		agg = reducer(agg, a.values[i])
	}

	if divBySize {
		return agg / float64(size)
	}
	return agg
}

func (a *Aggregator) Aggregate() float64 {
	switch a.aggType {
	case AggregatorType_MIN:
		return a.aggregateFn(reducerMin, false)
	case AggregatorType_MAX:
		return a.aggregateFn(reducerMax, false)
	default:
		// default is avg
		return a.aggregateFn(reducerSum, true)
	}
}
