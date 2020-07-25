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

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/net/websocket"
)

const DISTANCE_TOPIC_SENSOR_1 = "Distance.Sensor.1"
const DISTANCE_TOPIC_SENSOR_2 = "Distance.Sensor.2"

const BROKER = "tcp://localhost:1883"

var (
	addr = flag.String("addr", ":8080", "http service address")
)

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

	log.Println("+++ subscribeMQTT topic", topic)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		log.Printf("+++ Received %v", [2]string{msg.Topic(), string(msg.Payload())})
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

func distanceVSocket(ws *websocket.Conn) {
	distanceSocket(ws, DISTANCE_TOPIC_SENSOR_1, "1")
}

func distanceHSocket(ws *websocket.Conn) {
	distanceSocket(ws, DISTANCE_TOPIC_SENSOR_2, "2")
}

func distanceSocket(ws *websocket.Conn, topic, sensorNum string) {

	choke := make(chan [2]string) //[0] [1] - distance

	subscribeMQTT(choke, topic)

	sensor := "sensor" + sensorNum

	for {
		select {
		case incoming := <-choke:
			fmt.Printf("Sensor: %s\n", sensor)
			fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
			sendToWebSocket(ws, incoming[1])

			// color, err := convertDistanceToColor(incoming[1]) //publishMQTTLightOn(incoming[1])
			// if err != nil {
			// 	log.Println(err)
			// }
			// else {
			// 	publishMQTTLightOn(color)
			// 	sendToWebSocket(ws, incoming[1]+" "+color+" "+station) // distance color value e.g. "25 yellow"
			// }
		}
	}
}

// func publishMQTTWithPayload(topic, payload string) {
// 	opts := MQTT.NewClientOptions()
// 	opts.AddBroker(BROKER)
// 	opts.SetCleanSession(true)
// 	qos := 0

// 	client := MQTT.NewClient(opts)
// 	if token := client.Connect(); token.Wait() && token.Error() != nil {
// 		panic(token.Error())
// 	}

// 	token := client.Publish(topic, byte(qos), false, payload)
// 	token.Wait()
// 	if token.Error() != nil {
// 		fmt.Println(token.Error())
// 		os.Exit(1)
// 	}
// }
// func convertDistanceToColor(distanceData string) (string, error) {
// 	val, err := strconv.Atoi(distanceData)
// 	if err != nil {
// 		return "", err
// 	}

// 	switch {
// 	case val >= 0 && val <= 2:
// 		return "red-blink", nil
// 	case val > 2 && val <= 5:
// 		return "red", nil
// 	case val > 5 && val < 50:
// 		return "yellow", nil
// 	case val >= 50:
// 		return "green", nil
// 	}

// 	return "", fmt.Errorf("Value %s is out of range", distanceData)
// }
// func publishMQTTLightOn(data string) {
// 	publishMQTTWithPayload(LIGHT_ON_TOPIC_STATION_1, data)
// }

func sendToWebSocket(ws *websocket.Conn, data string) {

	//log.Printf("sendToWebSocket: \n%v\n", ws)

	websocket.Message.Send(ws, data)
}

func main() {
	http.HandleFunc("/", serveHome)
	http.Handle("/wshdistance", websocket.Handler(distanceHSocket))
	http.Handle("/wsvdistance", websocket.Handler(distanceVSocket))

	log.Printf("Started server %v\n", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
