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
	client, err := service.GetClientName()
	if err != nil {
		log.Println(err)
	}
	dto := rabbit.DTO{}
	dto.Data = client

	fmt.Println(client, err)
	err = b.Publish("New.Client", &dto)
}

func DirectoryScan(service config.Usecase, b *rabbit.Broker) {
	fmt.Println("reached")
	result, err := service.DirectoryScan("/")
	if err != nil {
		log.Println(err)
	}
	dto := rabbit.DTO{}
	dto.Data = result

	//fmt.Println(result, err)
	err = b.Publish("Directory.Scan", &dto)
}
