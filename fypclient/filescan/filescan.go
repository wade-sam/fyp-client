package filescan

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileData struct {
	Filepath string
	Filename string
	Checksum string
}
type FileScanResult struct {
	Filepath map[string]FileData
}

func DirectoryScan(startingPoint string, skip string) FileScanResult {
	tempHolder := make(map[string]FileData)
	os.Chdir(startingPoint)
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("prevent panic by handling failure accessing a path %q: %v\n", path, err)
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
			file := FileData{
				Filename: fileName,
				Filepath: filePath,
				Checksum: checkSum,
			}
			tempHolder[path] = file
		}

		//fmt.Println("vistited file or dir without errors %q\n", path)
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %v\n", err)
	}
	newFileScan := FileScanResult{Filepath: tempHolder}
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
