package backup

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"syscall"

	"github.com/wade-sam/fypclient/filescan"
)

func FullBackup(filescan filescan.FileScanResult) {
	for _, key := range filescan.Keys {
		value := filescan.Filepath[key]
		from, err := os.Open(key)
		if err != nil {
			log.Fatal(err)
		}
		defer from.Close()
		newFilePath := "/temp" + value.Filepath
		newFile := "/temp" + key
		//fmt.Println(key)
		CheckFolderExists(newFilePath, value.Filepath)
		CreateFileBackup(newFile, key, value.Permissions)
		//to, err := os.OpenFile(newFilePath)

	}
}

func CheckFolderExists(newFilePath, filePath string) {
	_, err := os.Stat(newFilePath)
	if os.IsNotExist(err) {
		CreateFolder(newFilePath, filePath)
	}

}

func CreateFolder(directory string, ogdirectory string) {
	m, _ := os.Stat(ogdirectory)
	perm := m.Mode().Perm()
	ogperms := fmt.Sprintf("%#o", perm)
	s, _ := strconv.ParseInt(ogperms, 0, 32)
	fmt.Println(s)
	os.MkdirAll(directory, os.FileMode(s))
	if own, ok := m.Sys().(*syscall.Stat_t); ok {
		os.Chown(directory, int(own.Uid), int(own.Gid))
	}
}

func CreateFileBackup(newfile string, existingFile string, filepermissions filescan.Permissions) {
	sourcefile, err := os.Open(existingFile)
	if err != nil {
		log.Fatal(err)
	}
	defer sourcefile.Close()
	newFile, err := os.Create(newfile)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(newFile, sourcefile); err != nil {
		fmt.Println("Error copying file", err)
	}
	newFile.Close()
	s, _ := strconv.ParseInt(filepermissions.Permissions, 0, 32)
	os.Chmod(newfile, os.FileMode(s))
	os.Chown(newfile, filepermissions.Ownership.UID, filepermissions.Ownership.GID)
	fmt.Println("Success", newfile)

}
