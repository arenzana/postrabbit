package main

import (
	"fmt"
	"log"

	"github.com/go-stomp/stomp"
)

func consume(config Config) {
	rabbitHost := fmt.Sprintf("%s:%s", config.RabbitHost, config.RabbitPort)
	conn, err := stomp.Dial("tcp", rabbitHost,
		stomp.ConnOpt.Login(config.RabbitUser, config.RabbitPassword),
		stomp.ConnOpt.AcceptVersion(stomp.V11),
		stomp.ConnOpt.AcceptVersion(stomp.V12),
		stomp.ConnOpt.Host(config.RabbitVHost),
		stomp.ConnOpt.Header("nonce", "B256B26D320A"))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Disconnect()
	subscription, err := conn.Subscribe(config.RabbitQueue, stomp.AckClient)
	if err != nil {
		log.Fatal(err)
	}

	for {
		payload := <-subscription.C
		log.Println("Received subscribed message -> ", string(payload.Body))
		err = conn.Ack(payload)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = subscription.Unsubscribe()
	if err != nil {
		log.Fatal(err)
	}
}
