package sensor

import (
	"fmt"
	"math/rand"
)

type Humidity struct {
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	UseRandom bool    `json:"useRandom"`
	State     bool    `json:"state"`
}

func (h Humidity) GetName() string {
	return h.Name
}

func (t Humidity) GetMQTTValue() string {
	mqttValue := t.GetValue()
	if t.UseRandom {
		mqttValue = t.GetRandomValue()
	}
	text := fmt.Sprintf("{%s: %f}", t.Name, mqttValue)
	return text
}

func (h Humidity) GetValue() float64 {
	return h.Value
}

func (h Humidity) GetState() bool {
	return h.State
}

func (h Humidity) GetRandomValue() float64 {
	min := 25.0
	max := 50.0
	r := min + r.Float64()*(max-min)
	return r
}
