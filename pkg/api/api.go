package api

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/fcproto/prototype/pkg/sensor"
)

type SensorData struct {
	Timestamp time.Time                `json:"timestamp"`
	Sensors   map[string]sensor.Values `json:"sensors"`
}

func NewSensorData() *SensorData {
	return &SensorData{
		Timestamp: time.Now(),
		Sensors:   make(map[string]sensor.Values),
	}
}

func (r *SensorData) String() string {
	var builder strings.Builder
	_ = json.NewEncoder(&builder).Encode(&r)
	return builder.String()
}
