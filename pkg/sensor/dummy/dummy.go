package dummy

import (
	"math/rand"
	"time"

	"github.com/fcproto/prototype/pkg/sensor"
)

type Dummy struct {
	rnd *rand.Rand
}

func (d *Dummy) Reset() {
	//TODO
}

func NewSensor() sensor.Sensor {
	return &Dummy{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (d *Dummy) GetValues() sensor.Values {
	temp := d.rnd.Float64() * 30
	return sensor.Values{
		"temp": temp,
	}
}
