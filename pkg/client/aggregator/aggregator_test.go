package aggregator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testAggType(t *testing.T, aggType Type, results []float64) {
	req := require.New(t)
	agg := New(aggType, 3)

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
	req.EqualValues(results[0], agg.Aggregate())

	// add 1 more value, this will replace the first value
	agg.AddValue(4)
	req.EqualValues(results[1], agg.Aggregate())
}

func TestAggregator(t *testing.T) {
	testCases := []struct {
		aggType Type
		results []float64
	}{
		{
			TypeAvg,
			[]float64{2, 3},
		},
		{
			TypeMin,
			[]float64{1, 2},
		},
		{
			TypeMax,
			[]float64{3, 4},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.aggType.String(), func(t *testing.T) {
			testAggType(t, testCase.aggType, testCase.results)
		})
	}
}
