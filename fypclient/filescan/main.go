package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type fileData struct {
	filepath string
	filename string
	checksum string
}
type fileScanResult struct {
	filepath map[string]fileData
}

//var files []fileData
//var fileStruct map[string]fileData
//var fileScanResult map[string]fileData

func main() {
	//filescanResult := make(map[string]fileData)
	fmt.Println("Hello you yute!")
	subDirToSkip := "golib"
	head := "/backup/Documents"
	//scanHolder := filesHolder
	fileScanResult := directoryScan(head, subDirToSkip)

	for key, value := range fileScanResult.filepath {
		fmt.Println(key, value.filename, value.checksum)
	}

	//	for _, value := range files {
	//	fmt.Println(value.filepath, value.filename, value.checksum)
	//}
}

func directoryScan(startingPoint string, skip string) fileScanResult {

	//fileOutput := make(map[string]fileData)
	tempHolder := make(map[string]fileData)

	filepath.Walk(startingPoint, func(path string, info os.FileInfo, err error) error {
		//outputScan := fileStruct
		if err != nil {
			//fmt.Println("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == skip {
			//fmt.Println("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}

		checkSum, err := hashFile(path)
		if err == nil {
			fileName := filepath.Base(path)
			filePath := filepath.Dir(path)
			file := fileData{
				filename: fileName,
				filepath: filePath,
				checksum: checkSum,
			}

			//fileScanResult[path] = file
			//outputScan[path] = file
			//files = append(files, file)
			tempHolder[path] = file
		}

		//fmt.Println("vistited file or dir without errors %q\n", path)
		return nil
	})
	newFileScan := fileScanResult{filepath: tempHolder}
	return newFileScan
}

func hashFile(filepath string) (string, error) {
	sha256sum := ""

	file, err := os.Open(filepath)
	if err != nil {
		return sha256sum, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return sha256sum, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnSha256String := hex.EncodeToString(hashInBytes)
	return returnSha256String, nil
}
