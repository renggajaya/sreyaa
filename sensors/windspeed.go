package sensor

import (
	"fmt"
	"math/rand"
)

type WindSpeed struct {
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	UseRandom bool    `json:"useRandom"`
	State     bool    `json:"state"`
}

func (h WindSpeed) GetName() string {
	return h.Name
}

func (t WindSpeed) GetMQTTValue() string {
	mqttValue := t.GetValue()
	if t.UseRandom {
		mqttValue = t.GetRandomValue()
	}
	text := fmt.Sprintf("{%s: %f}", t.Name, mqttValue)
	return text
}

func (h WindSpeed) GetValue() float64 {
	return h.Value
}

func (h WindSpeed) GetState() bool {
	return h.State
}

func (h WindSpeed) GetRandomValue() float64 {
	min := 0.0
	max := 100.0
	r := min + rand.Float64()*(max-min)
	return r
}
