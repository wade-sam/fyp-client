package rabbit

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type ConsumerConfig struct {
	ExchangeName string
	ExchangeType string
	RoutingKey   string
	QueueName    string
	ConsumerName string
	MaxAttempt   int
	Interval     time.Duration
	connection   *amqp.Connection
	Channels     *Channels
}

type Channels struct {
	Config  chan string
	Backup  chan DTO
	Restore chan DTO
}

type DTO struct {
	ID   string      `json:"id"`
	Data interface{} `json:"data"`
}

func (b *Broker) Start() (*amqp.Channel, error) {

	con, err := b.Connection()
	if err != nil {
		return nil, err
	}
	chn, err := con.Channel()
	if err != nil {
		return nil, err
	}
	b.channel = chn
	if err := chn.ExchangeDeclare(
		b.Consumer.ExchangeName,
		b.Consumer.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}
	if _, err := chn.QueueDeclare(
		b.Consumer.QueueName,
		true,
		false,
		false,
		false,
		amqp.Table{"x-message-ttl": 6000},
	); err != nil {
		return nil, err
	}

	if err := chn.QueueBind(
		b.Consumer.QueueName,
		b.Consumer.RoutingKey,
		b.Consumer.ExchangeName,
		false,
		nil,
	); err != nil {
		return nil, err
	}
	return chn, nil
}

func (b *Broker) Consume(channel *amqp.Channel) error {
	msgs, err := channel.Consume(
		b.Consumer.QueueName,
		b.Consumer.ConsumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for msg := range msgs {
		var d DTO
		//var backupDTO entity.ClientData

		//_, err := Deserialize(msg.Body)
		if err != nil {
			log.Println("Can't deserialise message", err)
		}

		switch msg.Type {

		case "Full.Backup":
			//dto := DTO{}
			//dto.Title = msg.Type
			//dto.Data = msg.Body
			//fmt.Println("DATA", d.Data)
			//err = json.Unmarshal([]byte(msg.Body), &d)
			//mapstructure.Decode(d.Data, &backupDTO)
			//backupDTO.Type = msg.Type
			// temp := entity.ClientData{
			// 	Type:       msg.Type,
			// 	Clientname: "samwade",
			// 	PolicyID:   "Friday Backup",
			// 	Data:       []string{},
			// }
			// d.Data = temp

			err = json.Unmarshal([]byte(msg.Body), &d)
			d.ID = msg.Type
			b.Consumer.Channels.Backup <- d
			fmt.Println("placed", d)
		case "restore":

		case "New.Client":
			fmt.Println("Consumer Received request")
			b.Consumer.Channels.Config <- msg.Type
		case "Directory.Scan":
			fmt.Println("Consumer Received request")
			b.Consumer.Channels.Config <- msg.Type
		}

		fmt.Println("msg consumed")
	}
	log.Println("Exiting")
	return nil
}
