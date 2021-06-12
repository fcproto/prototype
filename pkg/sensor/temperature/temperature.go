package temperature

import (
	"math"

	"github.com/fcproto/prototype/pkg/sensor"
)

type Temperature struct {
	t, baseTemp, amplitude float64
}

func (t *Temperature) Reset() {
	t.t = 0
}

func New(baseTemp, amplitude float64) sensor.Sensor {
	return &Temperature{
		t:         0,
		baseTemp:  baseTemp,
		amplitude: amplitude,
	}
}

func (t *Temperature) GetValues() sensor.Values {
	temp := t.baseTemp + t.amplitude*math.Sin(t.t)
	t.t += 0.05
	return sensor.Values{
		"temperature": temp,
	}
}
