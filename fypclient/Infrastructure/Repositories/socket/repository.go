package socket

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/wade-sam/fypclient/entity"
)

type Repository struct {
	Address    string
	Port       string
	Conn_Type  string
	Connection net.Conn
	Encoder    *gob.Encoder
}

type SockItem struct {
	ID     string      `json:"id"`
	Type   string      `json="type"`
	Client string      `json="client"`
	Item   interface{} `json="item"`
}

func NewRepository(address, port, conn_type string) *Repository {
	return &Repository{
		Address:   address,
		Port:      port,
		Conn_Type: conn_type,
	}
}

func (r *Repository) Connect() error {
	address := fmt.Sprintf("%v:%v", r.Address, r.Port)
	conn, err := net.Dial(r.Conn_Type, address)
	if err != nil {
		fmt.Println("connection error", err)
		return err
	}
	r.Connection = conn
	r.Encoder = gob.NewEncoder(conn)
	return nil

}

func (r *Repository) Disconnect() error {
	err := r.Connection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SendDirectoryLayout(item *SockItem) error {
	gob.Register(&entity.Directory{})
	item.Type = "directoryscan"
	err := r.Encoder.Encode(item)
	if err != nil {
		return err
	}
	fmt.Println("SEND DIRECTORY SCAN", item.ID)
	return nil
}

func (r *Repository) SendFile(item *SockItem) error {
	item.Type = "file"
	err := r.Encoder.Encode(item)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SendCompleteMessage(item *SockItem) error {
	item.Type = "clientcomplete"
	err := r.Encoder.Encode(item)
	if err != nil {
		return err
	}
	return nil

}
