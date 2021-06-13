package collector

import (
	"testing"

	"github.com/fcproto/prototype/pkg/client/aggregator"
	"github.com/fcproto/prototype/pkg/sensor"
	"github.com/stretchr/testify/require"
)

type mockSensor struct {
	data  []float64
	index int
}

func newMockSensor(data []float64) *mockSensor {
	return &mockSensor{
		data:  data,
		index: 0,
	}
}

func (m *mockSensor) Reset() {
	m.index = 0
}

func (m *mockSensor) GetValues() sensor.Values {
	ret := sensor.Values{
		"mock": m.data[m.index],
	}
	m.index++
	if m.index == len(m.data) {
		m.index = 0
	}
	return ret
}

func TestMockSensor(t *testing.T) {
	req := require.New(t)
	data := []float64{1, 2, 3}
	mSensor := newMockSensor(data)
	for i := 0; i < 10; i++ {
		for _, val := range data {
			req.Equal(val, mSensor.GetValues()["mock"])
		}
	}
}

func TestSensorCollector(t *testing.T) {
	req := require.New(t)
	data := []float64{1, 2, 3}
	sc := newSensorCollector(newMockSensor(data), []AggregateValues{{
		"mock": aggregator.TypeMax,
	}})
	req.EqualValues(0, sc.getValues()["mock"])
	for _, val := range data {
		sc.aggregate()
		req.Equal(val, sc.getValues()["mock"])
	}
	for i := 0; i < 10; i++ {
		sc.aggregate()
		req.EqualValues(3, sc.getValues()["mock"])
	}
}
