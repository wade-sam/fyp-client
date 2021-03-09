package rabbit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

type BrokerConfig struct {
	Schema         string
	Username       string
	Password       string
	Host           string
	Port           string
	VHost          string
	ConnectionName string
}

type Broker struct {
	config     BrokerConfig
	connection *amqp.Connection
	Producer   ProducerConfig
	Consumer   ConsumerConfig
}

func Deserialize(b []byte) (*DTO, error) {
	var msg DTO
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(msg)
	return &msg, err
}

func Serialize(msg *DTO) ([]byte, error) {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	err := encoder.Encode(msg)
	return b.Bytes(), err
}

func NewBroker(config BrokerConfig, producerconfig ProducerConfig, consumerconfig ConsumerConfig) *Broker {
	return &Broker{
		config:   config,
		Producer: producerconfig,
		Consumer: consumerconfig,
	}
}

func (r *Broker) Connect() error {
	if r.connection == nil || r.connection.IsClosed() {
		conn, err := amqp.Dial(fmt.Sprintf("%s://%s:%s@%s:%s/%s",
			r.config.Schema,
			r.config.Username,
			r.config.Password,
			r.config.Host,
			r.config.Port,
			r.config.VHost,
		))
		if err != nil {
			return err
		}
		r.connection = conn
	}
	return nil
}
func (r *Broker) Channel() (*amqp.Channel, error) {
	chn, err := r.connection.Channel()
	if err != nil {
		return nil, err
	}
	return chn, nil
}

func (r *Broker) Connection() (*amqp.Connection, error) {
	if r.connection == nil || r.connection.IsClosed() {
		return nil, errors.New("connection isnt open")
	}
	return r.connection, nil
}