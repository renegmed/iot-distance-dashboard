package mqtt_test

// import (
//         "sync"
//         "testing"

//         mqtt "github.com/eclipse/paho.mqtt.golang"
// )

// func TestMqttPubSub(t *testing.T) {

//         const TOPIC = "mytopic/test"

//         opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")

//         client := mqtt.NewClient(opts)
//         if token := client.Connect(); token.Wait() && token.Error() != nil {
//                 t.Fatal(token.Error())
//         }

//         var wg sync.WaitGroup
//         wg.Add(1)

//         if token := client.Subscribe(TOPIC, 0, func(client mqtt.Client, msg mqtt.Message) {
//    		defer wg.Done()
// 		if string(msg.Payload()) != "hello world!" {
//                         t.Fatalf("want 'hello world!', got %s", msg.Payload())
//                 }
//         }); token.Wait() && token.Error() != nil {
//                 t.Fatal(token.Error())
//         }

//         if token := client.Publish(TOPIC, 0, false, "hello world!"); token.Wait() && token.Error() != nil {
//                 t.Fatal(token.Error())
//         }
//         wg.Wait()
// }
