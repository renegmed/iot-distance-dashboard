package test

import (
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
	//opts.SetClientID(&ID)
	//opts.SetUsername(*user)
	//opts.SetPassword(*password)
	opts.SetCleanSession(CLEAN_SESSION)

	ch := make(chan [2]string)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		t.Logf("+++ Received %v", [2]string{msg.Topic(), string(msg.Payload())})
		ch <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	mqtt := cmd.ConfiguredMQTT{*opts, 1}

	htable := map[int]int{0: 70, 1: 60, 2: 50, 3: 40, 4: 30, 5: 20, 6: 10, 7: 5, 8: 0}
	vtable := map[int]int{9: 15, 10: 12, 11: 9, 12: 6, 13: 3, 14: 1}
	NUM := 14
	for i := 0; i <= NUM; i++ {

		if i < 9 {
			dist, _ := htable[i]
			//dist :=  70 - (i * 10)
			sDist := strconv.Itoa(dist)
			t.Log("Publishing distance", sDist)
			err := mqtt.Publish("Distance.Sensor.2", sDist)
			if err != nil {
				log.Println(err)
				t.Fail()
			}
		} else {
			dist, _ := vtable[i]

			sDist := strconv.Itoa(dist)
			t.Log("Publishing distance", sDist)
			err := mqtt.Publish("Distance.Sensor.1", sDist)
			if err != nil {
				log.Println(err)
				t.Fail()
			}
		}

		time.Sleep(1 * time.Second)
	}

	for i := 0; i < NUM; i++ {
		data, err := mqtt.Subscribe("Distance.Sensor.2", ch)
		if err != nil {
			log.Println(err)
			t.Fatal(err)
		}
		t.Log("Received data", data)
	}

}
