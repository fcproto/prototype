package speed

import (
	"github.com/fcproto/prototype/pkg/sensor"
)

type Speed struct {
	curr_speed   float64
	accelerating bool
}

func (s *Speed) Reset() {
	s.curr_speed = 0.0
	s.accelerating = true
}

func NewSensor() sensor.Sensor {
	return &Speed{
		curr_speed: 0.0,
	}
}

func (s *Speed) GetValues() sensor.Values {
	curr_speed := s.curr_speed
	if s.curr_speed > 2 {
		s.accelerating = false
	} else if s.curr_speed < 1 {
		s.accelerating = true
	}
	if s.accelerating {
		s.curr_speed += 1.0
	} else {
		s.curr_speed -= 1.0
	}
	return sensor.Values{
		"curr_speed": curr_speed,
	}
}
