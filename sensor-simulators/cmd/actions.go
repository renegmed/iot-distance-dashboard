package cmd

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type ConfiguredMQTT struct {
	Opts MQTT.ClientOptions
	Qos  int
}

func (c *ConfiguredMQTT) Publish(topic, payload string) error {
	client := MQTT.NewClient(&c.Opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	defer client.Disconnect(250)
	fmt.Println("Sample Publisher Started")
	// for i := 0; i < *num; i++ {
	// 	fmt.Println("---- doing publish ----")
	token := client.Publish(topic, byte(c.Qos), false, payload)
	token.Wait()
	// }

	fmt.Println("Sample Publisher Disconnected")
	return nil
}

func (c *ConfiguredMQTT) Subscribe(topic string, choke <-chan [2]string) (string, error) {
	//receiveCount := 0
	//choke := make(chan [2]string)

	// c.Opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
	// 	choke <- [2]string{msg.Topic(), string(msg.Payload())}
	// })

	client := MQTT.NewClient(&c.Opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return "", token.Error()
	}
	defer client.Disconnect(250)

	if token := client.Subscribe(topic, byte(c.Qos), nil); token.Wait() && token.Error() != nil {
		// fmt.Println(token.Error())
		// os.Exit(1)
		return "", token.Error()
	}

	// for receiveCount < *num {
	incoming := <-choke
	fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
	// 	receiveCount++
	// }

	fmt.Println("Sample Subscriber Disconnected")
	return incoming[1], nil

}
