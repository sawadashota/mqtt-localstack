// https://www.eclipse.org/paho/index.php?page=clients/golang/index.php
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	if len(os.Args) != 2 {
		log.Println(os.Args)
		log.Fatal("one argument required")
	}

	switch os.Args[1] {
	case "publisher":
		startPublisher()
	case "subscriber":
		startSubscriber()
	default:
		log.Fatal("argument should be publisher or subscriber")
	}
}

const topic = "sample/greeting"

func newClient(clientID string) mqtt.Client {
	ops := mqtt.NewClientOptions()
	ops.SetClientID(clientID)
	ops.AddBroker("tcp://broker:1883")
	return mqtt.NewClient(ops)
}

func startPublisher() {
	cl := newClient("publisher")

	if token := cl.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	for {
		select {
		case <-ctx.Done():
			cancel()
			log.Println("disconnecting...")
			cl.Disconnect(250)
			return
		default:
			token := cl.Publish(topic, 1, false, fmt.Sprintf("Hi, it's %s", time.Now().Format(time.RFC3339)))
			token.Wait()
			log.Println("published")
			time.Sleep(3 * time.Second)
		}
	}
}

func startSubscriber() {
	cl := newClient("subscriber")

	if token := cl.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	log.Println("starting subscriber")

	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("TOPIC: %s\n", msg.Topic())
		log.Printf("MSG: %s\n", msg.Payload())
	}

	if token := cl.Subscribe(topic, 1, messageHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	time.Sleep(30 * time.Second)

	if token := cl.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	log.Println("disconnecting...")
	cl.Disconnect(250)
}
