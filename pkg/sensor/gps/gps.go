package gps

import (
	"math"
	"time"

	"github.com/fcproto/prototype/pkg/sensor"
)

type GPS struct {
	factor      float64
	strech      float64
	seed        float64
	lat         float64
	lon         float64
	last_called float64
}

func (gps *GPS) Reset() {
}

func NewSensor() sensor.Sensor {
	return &GPS{
		strech:      5.,
		factor:      8.,
		seed:        0,
		lat:         52.514659,
		lon:         13.352144,
		last_called: float64(time.Now().Unix()),
	}
}

func (gps *GPS) GetValues() sensor.Values {
	curr_time := float64(time.Now().Unix())
	speed := (math.Sin(curr_time/gps.strech) + 1) * gps.factor
	speed_one_second_ago := (math.Sin((curr_time-1.)/gps.strech) + 1) * gps.factor
	acceleration := speed - speed_one_second_ago

	EARTH_RADIUS := 6378137.
	time_passed := curr_time - gps.last_called
	// avg_speed := gps.factor

	// we only move east :D
	dx := gps.factor * time_passed
	dy := 0.
	gps.lon += (180. / math.Pi) * (dx / EARTH_RADIUS) / math.Cos(math.Pi/180.*gps.lat)
	gps.lat += (180. / math.Pi) * (dy / EARTH_RADIUS)

	gps.last_called = curr_time

	return sensor.Values{
		"lat":          gps.lat,
		"lon":          gps.lon,
		"speed":        speed,
		"acceleration": acceleration,
	}
}
