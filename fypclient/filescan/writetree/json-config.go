package writetree

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/wade-sam/fypclient/filescan"
)

type File struct {
	FileName      string `json:"filename"`
	DirectoryPath string `json:"directory-path"`
	Permissions   string `json:"file-permissions"`
	UID           int    `json:"file-user-id"`
	GID           int    `json:"file-group-id"`
	Checksum      string `json:"checksum"`
}

type OutputFile struct {
	PolicyName string `json:"policy-name"`
	ScanDate   string `json:"scan-date"`
	//Files      struct {
	//	FileHolder map[string]File `file`
	//}
	Files map[string]File `files`
}

func WriteToFile(filescanresult filescan.FileScanResult) {
	outputfile := ObjectToJson(filescanresult)
	current_time := time.Now()

	filename := "/go/src/github.com/wade-sam/fypclient/file-config/Backup-" + current_time.Format("01-02-2006 15:04:05") + ".json"
	//os.Chdir("/go/src/github.com/wade-sam/fypclient/file-config")
	file, err := json.MarshalIndent(outputfile, "", "	")
	checkError(err)
	_ = ioutil.WriteFile(filename, file, 0775)
}

func ObjectToJson(filescanobject filescan.FileScanResult) OutputFile {
	tempHolder := make(map[string]File)
	for _, key := range filescanobject.Keys {
		value := filescanobject.Filepath[key]
		newfile := File{
			FileName:      value.Filename,
			DirectoryPath: value.Filepath,
			Permissions:   value.Permissions.Permissions,
			UID:           value.Permissions.Ownership.UID,
			GID:           value.Permissions.Ownership.GID,
			Checksum:      value.Checksum,
		}
		tempHolder[key] = newfile
	}
	outputfile := OutputFile{PolicyName: "policy1", ScanDate: "26/11/2020:15.13", Files: tempHolder}
	return outputfile
}

func CompareJsonFile(filescanresult filescan.FileScanResult) []string {
	newFiles := []string{}
	newscan := ObjectToJson(filescanresult)
	oldscan := ReadInJsonFile()

	for key, value := range newscan.Files {
		val, ok := oldscan.Files[key]
		if ok != false {
			if value.Checksum != val.Checksum {
				newFiles = append(newFiles, key)
			}
		} else {
			newFiles = append(newFiles, key)
		}
	}

	return newFiles
}
func ReadInJsonFile() OutputFile {
	file := filescan.InitialDirectoryScan("/go/src/github.com/wade-sam/fypclient/file-config", "")
	oldfilename := "/go/src/github.com/wade-sam/fypclient/file-config/"
	var previousScan OutputFile
	for _, value := range file.Filepath {
		oldfilename = oldfilename + value.Filename
	}
	oldfile, err := os.Open(oldfilename)
	defer oldfile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(oldfile)
	jsonParser.Decode(&previousScan)
	return previousScan
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
