package speed

import (
	"math"
	"time"

	"github.com/fcproto/prototype/pkg/sensor"
)

type Speed struct {
	factor float64
	strech float64
	seed   float64
}

func (s *Speed) Reset() {
}

func NewSensor() sensor.Sensor {
	return &Speed{
		strech: 5,
		factor: 8,
		seed:   0,
	}
}

func (s *Speed) GetValues() sensor.Values {
	curr_time := float64(time.Now().Unix())
	curr_speed := (math.Sin(curr_time/s.strech) + 1) * s.factor
	speed_one_second_ago := (math.Sin((curr_time-1.)/s.strech) + 1) * s.factor
	acceleration := curr_speed - speed_one_second_ago

	return sensor.Values{
		"speed":        curr_speed,
		"acceleration": acceleration,
	}
}
