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

func (r *SensorData) String() string {
	var builder strings.Builder
	_ = json.NewEncoder(&builder).Encode(&r)
	return builder.String()
}

// TODO
type Service struct {
}

func (s *Service) SubmitSensorData(data *SensorData) {

}
