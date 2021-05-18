package sensor

type Values map[string]float64

type Sensor interface {
	Reset()
	GetValues() Values
}
