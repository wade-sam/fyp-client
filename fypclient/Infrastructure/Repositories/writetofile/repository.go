package writetofile

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/wade-sam/fypclient/entity"
)

type FileRepo struct{}

type FileStruct struct {
	ClientName    string `json:"clientname"`
	BackupServer  string `json:"backupserver"`
	StorageNode   string `json:"storagenode"`
	RabbitDetails *RabbitConfig
}

type RabbitConfig struct {
	Schema         string `json:"schema"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Host           string `json:"host"`
	Port           string `json:"port"`
	VHost          string `json:"vhost"`
	ConnectionName string `json:"conname"`
}

func NewFileRepo() *FileRepo {
	return &FileRepo{}
}

func ReadInJsonFile() (*FileStruct, error) {
	var file FileStruct
	jsonFile, err := os.Open("/home/sam/Documents/fyp-client/fypclient/Infrastructure/Repositories/writetofile/config.json")
	if err != nil {
		return nil, entity.ErrFileNotFound
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, entity.ErrCouldNotUnMarshallJSON
	}
	json.Unmarshal(byteValue, &file)
	return &file, nil
}

func WriteJsonFile(file *FileStruct) error {
	outputFile, err := json.MarshalIndent(file, "", "	")
	if err != nil {
		return entity.ErrCouldNotMarshallJSON
	}
	err = ioutil.WriteFile("/home/sam/Documents/fyp-client/fypclient/Infrastructure/Repositories/writetofile/config.json", outputFile, 0775)
	if err != nil {
		return entity.ErrCouldNotWriteToFile
	}
	return nil
}

func (f *FileRepo) GetClientName() (string, error) {
	file, err := ReadInJsonFile()
	if err != nil {
		return "", err
	}
	client := file.ClientName
	return client, nil
}

func (f *FileRepo) GetStorageNode() (string, error) {
	file, err := ReadInJsonFile()
	if err != nil {
		return "", err
	}
	storagenode := file.StorageNode
	return storagenode, nil
}

func (f *FileRepo) GetRabbitDetails() (*RabbitConfig, error) {
	file, err := ReadInJsonFile()
	if err != nil {
		return nil, err
	}
	rabbit := file.RabbitDetails
	return rabbit, nil
}

func (f *FileRepo) SetClientName(name string) error {
	file, err := ReadInJsonFile()
	if err != nil {
		return err
	}
	file.ClientName = name
	err = WriteJsonFile(file)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileRepo) SetStorageNode(ip string) error {
	file, err := ReadInJsonFile()
	if err != nil {
		return err
	}
	file.StorageNode = ip
	err = WriteJsonFile(file)
	if err != nil {
		return err
	}
	return nil
}
