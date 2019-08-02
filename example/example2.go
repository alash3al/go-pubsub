package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/alash3al/go-pubsub"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	broker := pubsub.NewBroker()

	subscriber1, err := broker.Attach()
	if err != nil {
		log.Println("Error", err.Error())
	}

	broker.Subscribe(subscriber1, "BTCUSD")
	broker.Subscribe(subscriber1, "BTCAUD")
	broker.Subscribe(subscriber1, "BTCEUR")

	fmt.Println("Subscribers: ", broker.Subscribers("BTCUSD")) // Subscribers:  2

	fmt.Println(subscriber1.GetTopics()) // [BTCUSD]

	ch1 := subscriber1.GetMessages()

	go send("BTCUSD", broker)
	go send("BTCAUD", broker)
	go send("BTCEUR", broker)

	go receive(subscriber1.GetID(), ch1)

	fmt.Scanln()
	fmt.Println("done")
}

func getPrice() float64 {
	return rand.Float64() * 100000
}

func send(topic string, broker *pubsub.Broker) {
	fmt.Println("Sending...")
	for {
		p := getPrice()
		fmt.Printf("%v sending: %v\n", topic, p)
		broker.Broadcast(p, topic)
		time.Sleep(time.Second)
	}
}

func receive(id string, ch <-chan *pubsub.Message) {
	fmt.Printf("Subscriber %v, receiving...\n", id)
	for {
		if msg, ok := <-ch; ok {
			//fmt.Printf("Subscriber %v, received: %v\n", id, msg.GetPayload())
			fmt.Println(msg)
		}
	}
}
