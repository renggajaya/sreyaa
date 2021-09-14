package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"strings"
	"sync"

	"github.com/renggajaya/sreeya/sensor"
	"github.com/eclipse/paho.mqtt.golang"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Topic: %s\n", msg.Topic())
	fmt.Printf("Message: %s\n", msg.Payload())
}
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection Lost: %s\n", err.Error())
}


type RpcMessage struct {
	Method string `json:"method"`
	Params bool   `json:"params"`
}


func main() {
	var broker = "tcp://test.mosquitto.org:1883"
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("go_mqtt_sreeya")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectionLostHandler
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)


	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connected to server ", broker)

	if token := c.Subscribe("hardtmann/devices/sreeyaa/telemetry", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println("Error!")
		fmt.Println(token.Error())
		os.Exit(1)
	}

	var callback MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		topic := msg.Topic()
		payload := msg.Payload()
		s := string(payload[:])
		messageID := msg.MessageID()
		request := "hardtmann/devices/mesreeya/rpc/request/"
		requestID := topic[len(request):len(topic)]

		var payloadMessage RpcMessage
		if err := json.Unmarshal(payload, &payloadMessage); err != nil {
			panic(err)
		}

		fmt.Printf("TOPIC: %s\n", topic)
		fmt.Printf("MSG: %s\n", payload)
		fmt.Printf("payload: %s\n", s)
		fmt.Printf("MSG ID: %d\n", messageID)
		fmt.Printf("MSG Method: %s\n", payloadMessage.Method)
		fmt.Printf("MSG Method: %s\n", requestID)
		response, err := json.Marshal(payloadMessage)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		client.Publish("hardtmann/devices/sreeya/rpc/response/"+requestID, 0, false, string(response))

	}

	if token2 := c.Subscribe("hardtmann/devices/sreeya/rpc/request/+", 0, callback); token2.Wait() && token2.Error() != nil {
		fmt.Println("Error!")
		fmt.Println(token2.Error())
		os.Exit(1)
	}

	tempSensor := sensor.TemperatureSensor{Name: "temperature", Value: 0.0, UseRandom: true, State: true}
	humidSensor := sensor.HumiditySensor{Name: "humidity", Value: 0.0, UseRandom: true, State: true}
	WindSpeed := sensor.WindSpeed{Name: "wind speed", Value: 0.0, UseRandom: false, State: true}

	for {

		token := c.Publish("hardtmann/devices/sreeyaa/telemetry", 0, false, tempSensor.GetMQTTValue())
		token.Wait()

		token = c.Publish("hardtmann/devices/sreeyaa/telemetry", 0, false, humidSensor.GetMQTTValue())
		token.Wait()

		token = c.Publish("hardtmann/devices/sreeyaa/telemetry", 0, false, WindSpeed.GetMQTTValue())
		token.Wait()

		time.Sleep(1 * time.Second)

	}

	time.Sleep(6 * time.Second)

	if token := c.Unsubscribe("hardtmann/devices/sreeyaa/telemetry"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)

	time.Sleep(1 * time.Second)
}

	
	