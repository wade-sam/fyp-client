package handler

import (
	//"github.com/streadway/amqp"

	"fmt"
	"log"

	"github.com/wade-sam/fypclient/Infrastructure/Repositories/rabbit"
	config "github.com/wade-sam/fypclient/usecase/configuration"
)

func ConfigurationHandler(service config.Usecase, b *rabbit.Broker, chn chan string) {
	for msg := range chn {
		switch msg {
		case "New.Client":
			go GetClient(service, b)
		case "Directory.Scan":
			fmt.Println("reached")
			go DirectoryScan(service, b)
		}
	}
}

func GetClient(service config.Usecase, b *rabbit.Broker) {
	log.Println("trying to connect")
	client, err := service.ConfigureNewConsumerID()
	if err != nil {
		log.Println(err)
	}
	dto := rabbit.DTO{}
	dto.Data = client

	fmt.Println(client, err)
	err = b.Publish("New.Client", &dto)
	err = b.Disconnect()
	if err != nil {
		log.Println("could not disconnect from rabbit broker")
	}
	b.Consumer.RoutingKey = client
	b.Consumer.QueueName = client
	b.Consumer.ConsumerName = client

	err = b.Connect()
	if err != nil {
		log.Fatal("ERR", err)
	}
	consumer_chan, err := b.Start()
	if err != nil {
		log.Fatal("ERR", err)
	}
	go b.Consume(consumer_chan)
	log.Println("SUCCESS")

	// err = b.Connect()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// b.
	// consumer_chan, err := b.Start()
	// if err != nil {
	// 	log.Fatal("ERR", err)
	// }
	// log.Println(consumer_chan)
	// go b.Consume(consumer_chan)
}

func DirectoryScan(service config.Usecase, b *rabbit.Broker) {
	fmt.Println("reached")
	result, err := service.DirectoryScan("/")

	if err != nil {
		log.Println(err)
	}
	dto := rabbit.DTO{}
	dto.Data = result
	fmt.Println("sent", err)
	err = b.Publish("Directory.Scan", &dto)
}
