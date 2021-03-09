package handler_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wade-sam/fypclient/Infrastructure/Repositories/rabbit"
	"github.com/wade-sam/fypclient/Infrastructure/Repositories/writetofile"

	//"github.com/wade-sam/fypclient/pubsub/"
	"github.com/wade-sam/fypclient/pubsub/handler"
	"github.com/wade-sam/fypclient/usecase/configuration"
)

func initialise() (*rabbit.BrokerConfig, *rabbit.ProducerConfig, *rabbit.ConsumerConfig, *rabbit.Channels) {

	brokerconfig := rabbit.BrokerConfig{
		Schema:         "amqp",
		Username:       "admin",
		Password:       "85v!AP",
		Host:           "127.0.0.1",
		Port:           "5672",
		VHost:          "/",
		ConnectionName: "client1",
	}
	c := make(chan string)
	channs := rabbit.Channels{
		Config: c,
	}
	consumerconfig := rabbit.ConsumerConfig{
		ExchangeName: "main",
		ExchangeType: "direct",
		RoutingKey:   "host1",
		QueueName:    "host1",
		ConsumerName: "host1",
		MaxAttempt:   60,
		Interval:     1 * time.Second,
		Channels:     &channs,
	}

	producerConfig := rabbit.ProducerConfig{
		ExchangeName: "main",
		ExchangeType: "direct",
		MaxAttempt:   60,
		Interval:     1 * time.Second,
	}

	return &brokerconfig, &producerConfig, &consumerconfig, &channs
}
func Test_GetClientHandler(t *testing.T) {
	wtf := writetofile.NewFileRepo()
	brokerconfig, producerConfig, consumerconfig, channs := initialise()
	broker := rabbit.NewBroker(*brokerconfig, *producerConfig, *consumerconfig)
	err := broker.Connect()
	assert.Nil(t, err)
	service := configuration.NewConfigurationService(wtf)
	go handler.ConfigurationHandler(service, broker, channs.Config)
	for i := 1; i < 1; i++ {
		channs.Config <- "Directory.Scan"
	}

	select {}

}
