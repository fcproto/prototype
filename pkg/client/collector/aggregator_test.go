package collector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAvgAggregator(t *testing.T) {
	req := require.New(t)
	agg := NewAggregator(AggregatorType_AVG, 3)

	// aggregate should return 0 if no values have been added
	req.EqualValues(0, agg.Aggregate())

	// if only one value added, aggregate should return the single value
	agg.AddValue(99)
	req.EqualValues(99, agg.Aggregate())

	// fill with data
	testValues := []float64{1, 2, 3}
	for _, val := range testValues {
		agg.AddValue(val)
	}
	req.EqualValues(2, agg.Aggregate())

	// add 1 more value, this will replace the first value
	agg.AddValue(4)
	req.EqualValues(3, agg.Aggregate())
}
