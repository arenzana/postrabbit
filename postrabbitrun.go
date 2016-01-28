package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-stomp/stomp"
	pq "github.com/lib/pq"
)

func errorReporter(ev pq.ListenerEventType, err error) {
	if err != nil {
		log.Print(err)
	}
}

func run(config Config) {
	purl := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName, config.SslMode)
	listener := pq.NewListener(purl, 10*time.Second, time.Minute, errorReporter)
	err := listener.Listen("usertrigger")
	if err != nil {
		log.Fatal(err)
	}

	rabbitchannel := make(chan string, 100)

	//Code for STOMP
	go func() {
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

		for {
			payload := <-rabbitchannel
			log.Println(payload)
			err = conn.Send(config.RabbitQueue, "text/plain", []byte(payload))
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	for {
		select {
		case notification := <-listener.Notify:
			rabbitchannel <- notification.Extra
		case <-time.After(90 * time.Second):
			go func() {
				err := listener.Ping()
				if err != nil {
					log.Fatal(err)
				}
			}()
		}
	}
}
