package handler

import (
	"fmt"
	"log"
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/wade-sam/fypclient/Infrastructure/Repositories/rabbit"
	"github.com/wade-sam/fypclient/Infrastructure/Repositories/socket"
	"github.com/wade-sam/fypclient/entity"
	"github.com/wade-sam/fypclient/usecase/backup"
	"github.com/wade-sam/fypclient/usecase/configuration"
)

type ClientData struct {
	Type     string   `json:"type"`
	Client   string   `json:"name"`
	PolicyID string   `json:"policy"`
	Data     []string `json:"ignorelist"`
}

func BackupHandler(service backup.Usecase, configservice configuration.Usecase, b *rabbit.Broker, s *socket.Repository, chn chan rabbit.DTO) {
	for msg := range chn {
		switch msg.ID {
		case "Inc.Backup":
			ignore := []string{}
			mapstructure.Decode(msg, &ignore)

			err := StartIncrementalBackup(service, configservice, b, ignore)
			if err != nil {
				log.Println(err)
			}
		case "Full.Backup":
			var bdto ClientData
			//bdto := StoragenodeData{}
			err := mapstructure.Decode(msg.Data, &bdto)
			if err != nil {
				fmt.Println("ERROR", err)
			}
			fmt.Println("MAPSTRUCTURE", bdto)
			StartFullBackup(service, configservice, b, s, &bdto)
		}
	}
}

func StartIncrementalBackup(service backup.Usecase, configservice configuration.Usecase, b *rabbit.Broker, ignorelist []string) error {
	data, err := configservice.GetBackupResult()
	if err != nil {
		return err
	}
	fmt.Println(data)
	return nil
}

func StartFullBackup(service backup.Usecase, configservice configuration.Usecase, b *rabbit.Broker, s *socket.Repository, bdto *ClientData) {
	//write_files_to_disk := make(map[string]*entity.FileDTO)
	ignoremap := make(map[string]string)
	for i := range bdto.Data {
		ignoremap[bdto.Data[i]] = bdto.Data[i]
	}
	fmt.Println("REACHED")
	fmt.Println("map", ignoremap)
	head := "/home/sam/Documents"
	scanresult, err := service.BackupDirectoryScan(head, ignoremap)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("COMPLETED SCAN")
	SNtree := createStorageNodeTree(head, *scanresult)

	err = s.Connect()
	if err != nil {
		log.Println(err)
	}
	socket_directory := socket.SockItem{
		ID:     bdto.PolicyID,
		Client: bdto.Client,
		Item:   SNtree,
	}
	err = s.SendDirectoryLayout(&socket_directory)
	if err != nil {
		log.Println(err)
	}

	// workTree := createTotalTree("/", *scanresult)
	// backupchn := make(chan backup.FileTransfer)
	// producerchn := make(chan rabbit.DTO)
	// err = b.PublishManyInitialise(producerchn, "Client.File")
	// if err != nil {
	// 	log.Println("ERROR")
	// }
	// go service.FullBackupCopy(workTree, backupchn)
	// for msg := range backupchn {
	// 	dto := rabbit.DTO{
	// 		Title: "Client.File",
	// 		Data:  msg.BSFile,
	// 	}
	// 	producerchn <- dto
	// 	write_files_to_disk[msg.BSFile.ID] = msg.BSFile
	// 	//fmt.Println("sent", msg.BSFile.ID)
	// }
	// close(producerchn)
	// err = configservice.WriteBackupResult(write_files_to_disk)
	// if err != nil {
	// 	log.Println(err)
	// }
	// //bsscan, _, err := service.DirectoryScan("/", ignoremap)

	// sndto := rabbit.DTO{}
	// sndto.Data = &SNtree
	//err =   ("Client.Job", &dto)

}

func createStorageNodeTree(head string, scanresult map[string]interface{}) *entity.Directory {
	head = path.Clean(head)
	var snresult *entity.Directory
	for key, value := range scanresult {
		var parent *entity.Directory
		if key == head {
			fmt.Println("KEY: ", key, head, path.Base(key))
			snresult = value.(*entity.Directory)
			continue
		} else {
			parent = scanresult[path.Dir(key)].(*entity.Directory)
		}
		switch v := value.(type) {
		case *entity.Directory:
			parent.Folders[v.Name] = v //add value to the parents folders.
		}
	}

	return snresult
}
func createTotalTree(head string, scanresult map[string]interface{}) *entity.Directory {
	var snresult *entity.Directory
	for key, value := range scanresult {
		var parent *entity.Directory
		if key == head {
			snresult = value.(*entity.Directory)
			continue
		} else {
			parent = scanresult[path.Dir(key)].(*entity.Directory)
		}
		switch v := value.(type) {
		case *entity.Directory:
			parent.Folders[v.Name] = v //add value to the parents folders.
		case *entity.File: //If the interface type is *entity.File
			parent.Files = append(parent.Files, v) //append the file to the parent directory's File slice

		}
	}

	return snresult
}
