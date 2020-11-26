package writetree

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/wade-sam/fypclient/filescan"
)

type File struct {
	FileName      string `json:"filename"`
	DirectoryPath string `json:"directory-path"`
	Checksum      string `json:"checksum"`
}

type Files struct {
	Files map[string]File
}

type OutputFile struct {
	PolicyName string `json:"policy-name"`
	ScanDate   string `json:"scan-date"`
	Files      map[string]File
}

func WriteJSON(filescanresult filescan.FileScanResult) {
	tempHolder := make(map[string]File)

	for key, value := range filescanresult.Filepath {
		newfile := File{
			FileName:      value.Filename,
			DirectoryPath: value.Filepath,
			Checksum:      value.Checksum,
		}
		tempHolder[key] = newfile

	}
	outputfile := OutputFile{PolicyName: "policy1", ScanDate: "26/11/2020:15.13", Files: tempHolder}
	current_time := time.Now()

	filename := "Backup-" + current_time.Format("01-02-2006 15:04:05") + ".json"
	os.Chdir("/go/src/github.com/wade-sam/fypclient/file-config")
	file, err := json.MarshalIndent(outputfile, "", "	")
	checkError(err)
	_ = ioutil.WriteFile(filename, file, 0775)
}
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
