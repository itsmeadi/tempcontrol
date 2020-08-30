package models

type Reading struct {
	SensorID string  `json:"sensorID"`
	Type     string  `json:"type"`
	Value    float64 `json:"value"`
}

type Actuator struct {
	Level int `json:"level"`
}
