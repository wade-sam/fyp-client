package writetofile

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/wade-sam/fypclient/entity"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
<<<<<<< HEAD
	//jsonFile, err := os.Open("/home/sam/Documents/fyp-client/fypclient/Infrastructure/Repositories/writetofile/config.json")
	jsonFile, err := os.Open("config.json")
	if err != nil {
		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(path)
=======
	jsonFile, err := os.Open("/Users/sam/Documents/backup-client/fypclient/Infrastructure/Repositories/writetofile/config.json")
	if err != nil {
>>>>>>> 5134f6bfd5d5db16cc8a8142c4e964edd69fe139
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
	outputFile, err := json.MarshalIndent(file, "", " ")
	if err != nil {
		return entity.ErrCouldNotMarshallJSON
	}
<<<<<<< HEAD

	//err = ioutil.WriteFile("/home/sam/Documents/fyp-client/fypclient/Infrastructure/Repositories/writetofile/config.json", outputFile, 0775)
	err = ioutil.WriteFile("config.json", outputFile, 0775)
=======
	err = ioutil.WriteFile("/Users/sam/Documents/fyp-client/fypclient/Infrastructure/Repositories/writetofile/config.json", outputFile, 0775)
>>>>>>> 5134f6bfd5d5db16cc8a8142c4e964edd69fe139
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

func (f *FileRepo) CreateBackupResult(files map[string]*entity.FileDTO) error {
	output, err := json.MarshalIndent(files, "", "   ")
	if err != nil {
		fmt.Println(entity.ErrCouldNotMarshallJSON)
	}

	filename := "backup_config.gzip"
	file, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		fmt.Println(entity.ErrCouldNotMarshallJSON)
	}
	di, err := os.Create(filename)
	q := gzip.NewWriter(di)
	_, err = q.Write([]byte(file))
	if err != nil {
		return entity.ErrCouldNotWriteToFile
	}
	q.Close()
	return nil
}

func (f *FileRepo) GetPreviousBackupResult() (map[string]*entity.FileDTO, error) {
	var files map[string]*entity.FileDTO
	fi, err := os.Open("backup_config.gzip")
	if err != nil {
		return nil, err
	}
	reader, err := gzip.NewReader(fi)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(reader).Decode(&files)
	if err != nil {
		return nil, err
	}
	return files, nil
}
