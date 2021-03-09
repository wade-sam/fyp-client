package main

import (
	"log"
	"time"

	"github.com/wade-sam/fypclient/Infrastructure/Repositories/rabbit"
	"github.com/wade-sam/fypclient/Infrastructure/Repositories/writetofile"
	"github.com/wade-sam/fypclient/pubsub/handler"
	"github.com/wade-sam/fypclient/usecase/configuration"
)

func main() {
	wtf := writetofile.NewFileRepo()
	//conn_name, err := wtf.GetClientName()
	config_service := configuration.NewConfigurationService(wtf)

	conn_name, err := config_service.GetClientName()
	if err != nil {
		log.Fatal(err)
	}
	broker_config, err := wtf.GetRabbitDetails()
	if err != nil {
		log.Fatal(err)
	}
	BrokerConfig := rabbit.BrokerConfig{
		Schema:         broker_config.Schema,
		Username:       broker_config.Username,
		Password:       broker_config.Password,
		Host:           broker_config.Host,
		Port:           broker_config.Port,
		VHost:          broker_config.VHost,
		ConnectionName: conn_name,
	}
	channs := rabbit.Channels{
		Config:  make(chan string),
		Backup:  make(chan rabbit.DTO),
		Restore: make(chan rabbit.DTO),
	}
	consumerConf := rabbit.ConsumerConfig{
		ExchangeName: "main",
		ExchangeType: "direct",
		RoutingKey:   conn_name,
		QueueName:    conn_name,
		ConsumerName: conn_name,
		MaxAttempt:   60,
		Interval:     1 * time.Second,
		Channels:     &channs,
	}
	producerConf := rabbit.ProducerConfig{
		ExchangeName: "main",
		ExchangeType: "direct",
		MaxAttempt:   60,
		Interval:     1 * time.Second,
		RoutingKey:   "backupserver",
	}

	broker := rabbit.NewBroker(BrokerConfig, producerConf, consumerConf)
	err = broker.Connect()
	if err != nil {
		log.Fatal(err)
	}
	go handler.ConfigurationHandler(config_service, broker, channs.Config)
	consumer_chan, err := broker.Start()
	if err != nil {
		log.Fatal(err)
	}
	go broker.Consume(consumer_chan)

	select {}

}
