package writetree

import (
	"encoding/json"
	"fmt"

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
	FilesList  Files  `json:file-list`
}

func Writetest(filescanresult filescan.FileScanResult) {
	tempHolder := make(map[string]File)

	for key, value := range filescanresult.Filepath {
		newfile := File{
			FileName:      value.Filename,
			DirectoryPath: value.Filepath,
			Checksum:      value.Checksum,
		}
		tempHolder[key] = newfile

	}
	fmt.Println("Hello from the other side")
	//fmt.Println(tempHolder)
	files := Files{Files: tempHolder}
	outputfile := OutputFile{PolicyName: "policy1", ScanDate: "26/11/2020:15.13", FilesList: files}
	byteArray, err := json.MarshalIndent(outputfile, "", "	")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(byteArray))

}
