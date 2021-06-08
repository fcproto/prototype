package compass

import (
	"math/rand"
	"time"

	"github.com/fcproto/prototype/pkg/sensor"
)

type Compass struct {
	rnd      *rand.Rand
	rotation float64
}

func (c *Compass) Reset() {
	c.rotation = 0
}

func NewSensor() sensor.Sensor {
	return &Compass{
		rnd:      rand.New(rand.NewSource(time.Now().UnixNano())),
		rotation: 0,
	}
}

func (c *Compass) GetValues() sensor.Values {
	diff := rand.Float64() * 30
	switch c.rnd.Intn(2) {
	case 0:
		c.rotation += diff
	case 1:
		c.rotation -= diff
	}

	if c.rotation >= 360 {
		c.rotation -= 360
	}
	if c.rotation < 0 {
		c.rotation += 360
	}

	return sensor.Values{
		"rotation": c.rotation,
	}
}
