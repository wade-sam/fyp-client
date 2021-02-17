package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"time"
	"errors"
	"github.com/gorilla/mux"
	//"github.com/wade-sam/fypclient/backup"

	"github.com/wade-sam/fypclient/backup"
	"github.com/wade-sam/fypclient/filescan"
	"github.com/wade-sam/fypclient/filescan/writetree"
	"github.com/streadway/amqp"
)
type BrokerConfig struct{
	Schema         string
	Username       string
	Password       string
	Host           string
	Port           string
	VHost          string
	ConnectionName string
} 

type ConsumerConfig struct {
	ExchangeName  string
	ExchangeType  string
	RoutingKey    string
	QueueName     string
	ConsumerName  string
	Reconnect     struct {
		MaxAttempt int
		Interval   time.Duration
	}
}
type Broker struct {
	config BrokerConfig
	connection *amqp.Connection
}

func Filescan(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Filescan() called")
	w.Header().Set("Content-Type", "application/json")
	subDirToSkip := "golib"
	head := "/home/dev"
	filescan := filescan.InitialDirectoryScan(head, subDirToSkip)
	//fileScanResult := filescan.InitialDirectoryScan(head, subDirToSkip)
	response := writetree.ObjectToJson(filescan)

	json.NewEncoder(w).Encode(response)
	fmt.Println("Finished")
	//fmt.Fprintf(w, "Endpoint called:  test page")
}

func FBackup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("FullBackup called")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Full backup started")
	subDirToSkip := ""
	head := "/home/dev"
	filescan := filescan.InitialDirectoryScan(head, subDirToSkip)
	backup.FullBackup(filescan)
	fmt.Println("FullBackup called")
}

func IBackup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Incremental backup started")
	subDirToSkip := "golib"
	head := "/home/dev"
	filescan := filescan.InitialDirectoryScan(head, subDirToSkip)
	backup.IncrementalBackup(filescan)
	fmt.Println("Incremental() started")
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/filescan", Filescan).Methods("GET")
	router.HandleFunc("/full", FBackup).Methods("GET")
	router.HandleFunc("/incremental", IBackup).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}



func RabitCreateConnection() amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@192.168.1.210:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	//defer conn.Close()
	return *conn
}

func RabbitConnectToDefaultChannel(conns amqp.Connection) amqp.Channel {
	ch, err := conns.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return *ch

}

//func RabbitConsumeDefaultChannel(channel amqp.Channel){

//}

//NewBroker returns a RabbitMQ instance
func NewBroker(config BrokerConfig) *Broker{
	return &Broker{
		config: config,
	}
}

// Connection returns exiting `*amqp.Connection` instance.
func(r *Broker) Connection() (*amqp.Connection, error){
	if r.connection == nil || r.connection.IsClosed(){
		return nil, errors.New("connection isnt open")
	}
	return r.connection, nil
}

//Channel returns a new *amqp.Channel instance

func (r *Broker) Channel() (*amqp.Channel, error){
	chn, err := r.connection.Channel()
	if err != nil{
		return nil, err
	}
	return chn, nil
}

//Connects to the RabbitMQ server
func (r *Broker) Connect() error{
	if r.connection == nil || r.connection.IsClosed() {
		conn, err := amqp.Dial(fmt.Sprintf("%s://%s:%s@%s:%s/%s",
			r.config.Schema,
			r.config.Username,
			r.config.Password,
			r.config.Host,
			r.config.Port,
			r.config.VHost,
		))
		if err != nil{
			return err
		}
		r.connection = conn
		
	}
	return nil
}

type Consumer struct {
	config ConsumerConfig
	Broker *Broker
}

func NewConsumer(config ConsumerConfig, broker *Broker) *Consumer {
	return &Consumer{
		config: config,
		Broker: broker,
	}
}

func (c *Consumer) Start() error{
	con, err := c.Broker.Connection()
	if err != nil{
		return err
	}

	chn, err := con.Channel()
	if err != nil{
		return err
	}

	if err := chn.ExchangeDeclare(
		c.config.ExchangeName,
		c.config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil{
		return nil
	}
	if _, err:= chn.QueueDeclare(
		c.config.QueueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil{
		return nil
	}

	if err := chn.QueueBind(
		c.config.QueueName,
		c.config.RoutingKey,
		c.config.ExchangeName,
		false,
		nil,
	); err != nil {
		return err
	}
	go c.consume(chn)
	return nil
}

func (c *Consumer) consume(channel *amqp.Channel){
	msgs, err := channel.Consume(
		c.config.QueueName,
		c.config.ConsumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil{
		log.Println("Could not start consumer",err)
		return
	}

	for msg := range msgs{
		log.Println("Consumed", string(msg.Body))
		if err := msg.Ack(false); err != nil {
			log.Println("unable to acknowledge the message, dropped", err)
		}
	}
	
	log.Println("Exiting")

}


func main() {

	bc := BrokerConfig{
		Schema: "amqp",
		Username: "admin",
		Password: "85v!AP",
		Host: "192.168.1.210",
		Port: "5672",
	}

	consumerconfig := ConsumerConfig{
		ExchangeName: "amq.topic",
		ExchangeType: "topic",
		RoutingKey: "host1",
		QueueName: "host1",
		ConsumerName: "host1",
	}
	consumerconfig.Reconnect.MaxAttempt = 60
	consumerconfig.Reconnect.Interval = 1 * time.Second


	brokerInstance := NewBroker(bc)
	if err := brokerInstance.Connect(); err != nil{
		log.Fatal("Unable to connect to broker")
	}

	consumerInstance := NewConsumer(consumerconfig, brokerInstance)
	if err := consumerInstance.Start(); err != nil{
		log.Fatalln("Unable to start consumer", err)
	}
	select {}



	//conn := RabitConnectToConnection()
	//ch := RabbitCreateDefaultChannel(conn)

	//handleRequests()
	//subDirToSkip := "golib"
	//head := "/backup/Documents"
	//fileScanResult := filescan.InitialDirectoryScan(head, subDirToSkip)
	//for key, value := range fileScanResult.Filepath {
	//	fmt.Println(key, value.Filename, value.Checksum, value.Permissions.Ownership)
	//	}
	//backup.FullBackup(fileScanResult)
	//backup.IncrementalBackup(fileScanResult)
	//writetree.WriteToFile(fileScanResult)
	//time.Sleep(20 * time.Second)
	//	fmt.Println("Checking for differences")
	//differences := writetree.CompareJsonFile(fileScanResult)
	//fmt.Println(differences)
	//readfile := writetree.ReadInJsonFile()

	//fileScanResult := filescanDirectoryScan(head, subDirToSkip)
	//fmt.Println(fileScanResult)
	//for key, value := range fileScanResult.Filepath {
	//	fmt.Println(key, value.Filename, value.Checksum)
	//}

}
