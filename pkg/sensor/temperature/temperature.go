package temperature

import (
	"math"

	"github.com/fcproto/prototype/pkg/sensor"
)

type Temperature struct {
	key      string
	t        float64
	baseTemp float64
}

func (t *Temperature) Reset() {
	t.t = 0
}

func NewEnvironmentSensor() sensor.Sensor {
	return &Temperature{
		key:      "environment",
		t:        0,
		baseTemp: 25,
	}
}

func NewTrackSensor() sensor.Sensor {
	return &Temperature{
		key:      "track",
		t:        0,
		baseTemp: 30,
	}
}

func (t *Temperature) GetValues() sensor.Values {
	temp := t.baseTemp + 3*math.Sin(t.t)
	t.t += 0.05
	return sensor.Values{
		t.key: temp,
	}
}
