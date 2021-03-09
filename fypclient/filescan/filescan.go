package filescan

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

type FileData struct {
	Filepath    string
	Filename    string
	Permissions Permissions
	Checksum    string
}

type FileScanResult struct {
	Filepath map[string]FileData
	Keys     []string
	//Filepath *orderedmap.OrderedMap
}
type Ownership struct {
	UID int
	GID int
}
type Permissions struct {
	Ownership   Ownership
	Permissions string
}

func InitialDirectoryScan(startingPoint string, skip string) FileScanResult {

	tempHolder := make(map[string]FileData)
	keylist := []string{}
	os.Chdir(startingPoint)
	err := filepath.Walk(startingPoint, func(path string, info os.FileInfo, err error) error {
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
			permissions, ownership := GetPermissions(info)
			//fmt.Print(permissions, filepath.Dir(path))
			permHolder := Permissions{
				Ownership:   ownership,
				Permissions: permissions,
			}
			file := FileData{
				Filename:    fileName,
				Filepath:    filePath,
				Permissions: permHolder,
				Checksum:    checkSum,
			}
			tempHolder[path] = file
			keylist = append(keylist, path)

		}

		//fmt.Println("vistited file or dir without errors %q\n", path)
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %v\n", err)
	}
	newFileScan := FileScanResult{
		Filepath: tempHolder,
		Keys:     keylist,
	}
	//newFileScan := FileScanResult{Filepath: FilescanHolder}
	return newFileScan
}

func GetPermissions(info os.FileInfo) (string, Ownership) {
	Ownership := Ownership{}
	perm := info.Mode().Perm()
	if own, ok := info.Sys().(*syscall.Stat_t); ok {
		Ownership.UID = int(own.Uid)
		Ownership.GID = int(own.Gid)

	}
	permissions := fmt.Sprintf("%#o", perm)

	return permissions, Ownership

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
