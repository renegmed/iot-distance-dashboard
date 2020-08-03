package test

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/renegmed/iot-distance-dashboard/sensor-simulators/cmd"
)

const (
	CLIENT_ID = 1001
	BROKER    = "tcp://localhost:1883" // "tcp://iot.eclipse.org:1883"
	// USER          = ""
	// PASSWORD      = ""
	CLEAN_SESSION = true
)

func TestPublishSubscribe(t *testing.T) {

	t.Log("Start testing....")

	opts := MQTT.NewClientOptions()
	opts.AddBroker(BROKER)
	opts.SetCleanSession(CLEAN_SESSION)

	ch := make(chan [2]string)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		t.Logf("+++ Received %v", [2]string{msg.Topic(), string(msg.Payload())})
		ch <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	mqtt := cmd.ConfiguredMQTT{*opts, 1}

	// array of {horizontal, vertical} sensor values
	sensorsData := [][]int{{70, 0}, {60, 0}, {50, 0}, {40, 0}, {30, 0}, {20, 0}, {10, 0}, {5, 0}, {0, 15}, {0, 12}, {0, 9}, {0, 6}, {0, 3}, {0, 1}}
	NUM := len(sensorsData)

	for i := 0; i < NUM; i++ {
		dist2 := sensorsData[i][1]
		dist1 := sensorsData[i][0]

		t.Logf(fmt.Sprintf("Publishing sensors distance%d %d", dist1, dist2))
		err := mqtt.Publish("Distance.Sensors", strconv.Itoa(dist1)+" "+strconv.Itoa(dist2))
		if err != nil {
			log.Println(err)
			t.Fail()
		}

		time.Sleep(500 * time.Millisecond)
	}

	// for i := 0; i < NUM; i++ {
	// 	data, err := mqtt.Subscribe("Distance.Sensor.2", ch)
	// 	if err != nil {
	// 		log.Println(err)
	// 		t.Fatal(err)
	// 	}
	// 	t.Log("Received data", data)
	// }

}
