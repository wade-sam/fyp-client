package rabbit

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type ProducerConfig struct {
	ExchangeName string
	ExchangeType string
	RoutingKey   string
	MaxAttempt   int
	Interval     time.Duration
	connection   *amqp.Connection
}

func (b *Broker) Publish(Type string, body *DTO) error {
	channel, err := b.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	data, err := json.Marshal(body.Data)
	//dto, err := Serialize(body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(Type, "", b.Producer.RoutingKey)
	err = channel.Publish(
		b.Producer.ExchangeName,
		b.Producer.RoutingKey,
		false,
		false,
		amqp.Publishing{
			Type:        Type,
			ContentType: "encoding/json",
			Body:        []byte(data),
		},
	)
	if err != nil {
		return err
	}
	fmt.Println("Sent message back")
	return nil
}

func (b *Broker) PublishManyInitialise(c chan DTO, Type string) error {
	go Publishmany(c, b, Type)
	return nil

}

func Publishmany(c chan DTO, b *Broker, Type string) error {
	channel, err := b.Channel()
	if err != nil {
		log.Println(err)
		return err
	}
	defer channel.Close()
	for msg := range c {
		data, err := json.Marshal(msg.Data)
		if err != nil {
			log.Println(err)
		}
		if err := channel.Publish(
			b.Producer.ExchangeName,
			b.Producer.RoutingKey,
			false,
			false,
			amqp.Publishing{
				Type:        Type,
				ContentType: "encoding/json",
				Body:        []byte(data),
			},
		); err != nil {
			log.Println(err)
		}
	}
	return nil
}
