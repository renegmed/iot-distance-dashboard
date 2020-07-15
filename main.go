package main

/*

This simple web application shows how to communicate with sensors and LEDS hardwares
through MQTT broker.

This app:

	1. Receives distance data from distance sensor
	2. Pushes the distance data to a web page using websocket
	3. Based on the distance, this app determines what color of light to turn on
	4. Using MQTT broker, the app sends and publishes message to LED lights which color to turn on
    5. Finally, the app displays, on the web page ,the turned-on light.
*/
import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/net/websocket"
)

const DISTANCE_TOPIC = "Distance.Station.1"
const LIGHT_ON_TOPIC = "Light.On.Station.1"
const BROKER = "tcp://localhost:1883"

var (
	addr = flag.String("addr", ":8080", "http service address")
)
var distanceData string

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func subscribeMQTT(choke chan [2]string, topic string) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(BROKER)
	opts.SetCleanSession(true)
	qos := 0
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(topic, byte(qos), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

}

func distanceSocket(ws *websocket.Conn) {

	choke := make(chan [2]string)

	subscribeMQTT(choke, DISTANCE_TOPIC)

	for {
		select {
		case incoming := <-choke:
			log.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
			color, err := publishMQTTLightOn(incoming[1])
			if err != nil {
				log.Println(err)
			} else {
				sendToWebSocket(ws, incoming[1]+" "+color) // distance color value e.g. "25 yellow"
			}
		}
	}
}

func publishMQTTWithPayload(topic, payload string) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(BROKER)
	opts.SetCleanSession(true)
	qos := 0

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := client.Publish(topic, byte(qos), false, payload)
	token.Wait()
	if token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func publishMQTTLightOn(data string) (string, error) {
	val, err := strconv.Atoi(data)
	if err != nil {
		return "", err
	}
	var color string
	switch {
	case val >= 0 && val < 20:
		color = "red"
	case val >= 20 && val < 40:
		color = "yellow"
	case val >= 40:
		color = "green"
	}

	publishMQTTWithPayload(LIGHT_ON_TOPIC, color)

	return color, nil
}

func lightOnSocket(ws *websocket.Conn) {
	choke := make(chan [2]string)

	subscribeMQTT(choke, LIGHT_ON_TOPIC)

	for {
		select {
		case incoming := <-choke:
			log.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
			sendToWebSocket(ws, incoming[1])
		}
	}
}

func sendToWebSocket(ws *websocket.Conn, data string) {
	websocket.Message.Send(ws, data)
}

func main() {
	http.HandleFunc("/", serveHome)
	http.Handle("/wsdistance", websocket.Handler(distanceSocket))
	http.Handle("/wslighton", websocket.Handler(lightOnSocket))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
