package temperature

import (
	"math"

	"github.com/fcproto/prototype/pkg/sensor"
)

type Temperature struct {
	t        float64
	baseline float64
}

func (t *Temperature) Reset() {
	t.t = 0
}

func NewSensor() sensor.Sensor {
	return &Temperature{
		t:        0,
		baseline: 25,
	}
}

func (t *Temperature) GetValues() sensor.Values {
	temp := t.baseline + 3*math.Sin(t.t)
	t.t += 0.05
	return sensor.Values{
		"temp": temp,
	}
}
