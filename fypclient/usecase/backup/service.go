package backup

import (
	"container/list"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/wade-sam/fypclient/entity"
)

type Service struct {
	Files *entity.Directory
}

func NewBackupService() *Service {
	return &Service{
		Files: &entity.Directory{},
	}
}

// func (s *Service) StartIncrementalBackup(start, id string, ignorelist []string) error {
// 	ignoremap := make(map[string]string)
// 	for i := range ignorelist {
// 		ignoremap[ignorelist[i]] = ignorelist[i]
// 	}
// 	scan, tree, err := s.DirectoryScan(start, ignoremap)
// 	if err != nil {
// 		return err
// 	}
// 	s.bsrepo.ExpectedFiles(scan)
// 	s.snrepo.BackupLayout(tree)

// 	return nil
// }
type FileTransfer struct {
	Status string
	BSFile *entity.FileDTO
	SNFile *entity.File
	Data   *[]byte
}

/*
Full Backup Process
- Do an initial scan of the FS, just recording the directories, comparing against the ignore list. Store these in a map like structure.
- Caller function will then create a tree out of this and send it to the storage node. Pass this tree to the FullBackupCopy function, as well as a channel
FullBackupCopy function
	- Take the tree structure and a channel as params
	- Create a map structure which will hold the backed up files/directories
	- Map structure will be used to hold all the files in the FS, and store it to disk
	- Then walk through the tree with either a DFS/BFS algo. For each file add the metadata to the map structure, whilst also passing the metadata and data back on the channel using the FileTransfer struct.
	- Once done write the map sruct to file

Inc Backup Process
- Read in the most recent backup scan
- Do an initial scan of the FS, just recording the directories, comparing against the ignore list. Store in a map structure
Caller function will then create a tree out of this. Pass this tree to the IncBackupCopy function, as well as a channel
IncBackupCopy function(map, previousbackupscan, channel)
- Take the tree structure, previousbackupscan and a channel as params
- Create a map structure which will hold all the files the files in the file system
- Walk through the tree struct comparing items against those found in the previousbackupscan. If there are differences add them to the tree struct, as well as the map struct. If there are no differences only add to the map struct

Inc Backup Read in the most recent scan and walk through the file tree as usual, but compare to the most recent scan. If there are any changes record them, if there aren't don't


Steps Full:
	- Perform Directory scan of FS, using ignorelist. Store in a tree. Send to SN
	- Loop through tree and backup each file. Store metadata in a map struct
	- Once complete write map struct to file

Steps Inc:
	- Read in previous backup scan
	- Perform scan of FS, using ignorelist.
	- Compare newscan against old scan and store in a map struct. Extract directories from map struct and send to SN. Create a tree struct from map struct, which will hold new files and directories to backup.
	- Using tree struct walk through and backup each individual file

*/

func (s *Service) BackupDirectoryScan(head string, ignorelist map[string]string) (*map[string]interface{}, error) {
	head = path.Clean(head) //cleans up the path, returning the shortest filepath
	fmt.Println(head)
	var nodes = map[string]interface{}{}
	var walkfunc filepath.WalkFunc = func(loc string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
		} else if info.IsDir() {
			//fmt.Println("DIRECTORY", loc)
			if _, ok := ignorelist[loc]; ok {
				return filepath.SkipDir
			} else {
				n := entity.NewDirectory(path.Base(loc))
				n.Path = loc
				owner, uid, gid := GetProperties(info)
				n.AddProperties(owner, uid, gid)
				nodes[loc] = n

			}

		} else if info.Mode().IsRegular() {
			if _, err := os.Stat(loc); err == nil || os.IsExist(err) {

				SNFile := entity.NewFile(path.Base(loc))
				BSFile, err := entity.NewFileDTO(loc)
				if err != nil {
					return err
				}
				n := entity.NewFile(path.Base(loc))
				n.Path = loc
				SNFile.Path = loc
				owner, uid, gid := GetProperties(info)
				n.AddProperties(owner, uid, gid)
				SNFile.AddProperties(owner, uid, gid)
				checksum, err := md5checksum(loc)
				if err != nil {
					return err
				}
				n.AddChecksum(checksum)
				nodes[loc] = n
				err = BSFile.AddChecksum(checksum)

				SNFile.AddChecksum(checksum)

			} else {
				return err
			}
		}
		return nil
	}

	err := filepath.Walk(head, walkfunc)
	if err != nil {
		log.Println("Error Could not execute filepath walk", err)
		return nil, err

		//return nil, nil, err
	}

	return &nodes, nil
}

func (s *Service) IncBackupComparison() {

}

func (s *Service) FullBackupCopy(n *entity.Directory, chn chan (FileTransfer)) {
	visited := make(map[string]*entity.Directory)
	queue := list.New()
	queue.PushBack(n)
	visited[n.Path] = n
	for queue.Len() > 0 {
		pop := queue.Front()
		for _, node := range pop.Value.(*entity.Directory).Folders {
			if _, ok := visited[node.Path]; !ok {
				visited[node.Path] = node
				queue.PushBack(node)
			}
		}
		for _, value := range pop.Value.(*entity.Directory).Files {
			data, err := copy(value.Path)
			if err != nil {
				msg := entity.FileDTO{
					ID:       value.Path,
					Status:   "Failed",
					Checksum: value.Checksum,
				}
				fmt.Println("SENT FILE TYPE: Failed")
				dto := FileTransfer{
					Status: "Failed",
					BSFile: &msg,
				}
				chn <- dto
			} else {
				bsmsg := entity.FileDTO{
					ID:       value.Path,
					Status:   "Success",
					Checksum: value.Checksum,
				}
				dto := FileTransfer{
					Status: "Success",
					BSFile: &bsmsg,
					SNFile: value,
					Data:   data,
				}

				chn <- dto

			}
		}
		queue.Remove(pop)
	}
	fmt.Println("FINISHED COPYING FILES")
	final := entity.FileDTO{
		Status: "Finished",
	}
	dto := FileTransfer{
		Status: "Finished",
		BSFile: &final,
	}
	chn <- dto
	close(chn)
}

func copy(path string) (*[]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	filebytes, err := ioutil.ReadAll(file)
	//temp := []byte(file)
	if err != nil {
		return nil, err
	}
	return &filebytes, nil
}

// func (s *Service) FullBackupFiles(chn chan (FileTransfer), ignore map[string]string) (*entity.FileDTO, error) {
// 	var nodes = map[string]interface{}{}
// 	var walkfunc filepath.WalkFunc = func(loc string, info fs.FileInfo, err error) error {
// 		if err != nil {
// 			log.Println(err)
// 		} else if info.IsDir() {
// 			if len(ignore) > 0 {
// 				if _, ok := ignore[loc]; ok {
// 					return filepath.SkipDir
// 				}
// 			} else {
// 				n := entity.NewDirectory(path.Base(loc))
// 				owner, uid, gid := GetProperties(info)
// 				n.AddProperties(owner, uid, gid)

// 				nodes[loc] = n
// 			}
// 		} else if info.Mode().IsRegular() {
// 			if _, err := os.Stat(loc); err == nil || os.IsExist(err) {

// 				SNFile := entity.NewFile(path.Base(loc))
// 				BSFile, err := entity.NewFileDTO(loc)
// 				if err != nil {
// 					return err
// 				}
// 				n := entity.NewFile(path.Base(loc))
// 				n.Path = loc
// 				SNFile.Path = loc
// 				owner, uid, gid := GetProperties(info)
// 				n.AddProperties(owner, uid, gid)
// 				SNFile.AddProperties(owner, uid, gid)
// 				checksum, err := md5checksum(loc)
// 				if err != nil {
// 					return err
// 				}
// 				n.AddChecksum(checksum)
// 				nodes[loc] = n
// 				err = BSFile.AddChecksum(checksum)

// 				SNFile.AddChecksum(checksum)

// 			} else {
// 				return err
// 			}
// 		}

// 		return nil
// 	}
// 	err = filepath.Walk("/", walkfunc)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	close(chn)
// }

func (s *Service) DirectoryScan(head string, ignore map[string]string) (map[string]*entity.File, *entity.Directory, error) {
	fmt.Println("reached")
	head = path.Clean(head)
	//var sntree *entity.Directory
	bsfiles := make(map[string]*entity.File)
	var nodes = map[string]interface{}{}
	var walkfunc filepath.WalkFunc = func(loc string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
		} else if info.IsDir() {
			if len(ignore) > 0 {
				if _, ok := ignore[loc]; ok {
					return filepath.SkipDir
				}

			} else {
				n := entity.NewDirectory(path.Base(loc))
				owner, uid, gid := GetProperties(info)
				n.AddProperties(owner, uid, gid)

				nodes[loc] = n

			}
		} else if info.Mode().IsRegular() {
			if _, err := os.Stat(loc); err == nil || os.IsExist(err) {
				n := entity.NewFile(path.Base(loc))
				n.Path = loc
				owner, uid, gid := GetProperties(info)
				n.AddProperties(owner, uid, gid)
				checksum, err := md5checksum(loc)
				if err != nil {
					log.Println(err)
				}
				n.AddChecksum(checksum)
				nodes[loc] = n
				bsfiles[n.Path] = n

			} else {
				return err
			}
		}

		return nil
	}
	err := filepath.Walk(head, walkfunc)
	if err != nil {
		return nil, nil, err
	}
	// this keeps looping through until all the directories have been placed inside eachother.
	for key, value := range nodes {
		var parentDirectory *entity.Directory // initialise parentDirectory as a pointer of entity.Directory. Doesn't have an actual value yet

		if key == head {
			s.Files = value.(*entity.Directory) //if key is equal to head then assign value to S.Files
			continue
		} else {
			parentDirectory = nodes[path.Dir(key)].(*entity.Directory) //rooparentDirectory is NOW assigned as the parent directory of values directory. path.Dir takes off the last element in the path
		}
		switch v := value.(type) {
		case *entity.Directory: //If the interface type is *entity.Directory
			parentDirectory.Folders[v.Name] = v //add value to the parents folders.
		case *entity.File: //If the interface type is *entity.File
			parentDirectory.Files = append(parentDirectory.Files, v) //append the file to the parent directory's File slice
		}
	}

	return bsfiles, s.Files, nil
}

func GetProperties(info os.FileInfo) (string, string, string) {
	var uid string
	var gid string
	perm := info.Mode().Perm()
	owner := fmt.Sprintf("%#o", perm)
	if own, ok := info.Sys().(*syscall.Stat_t); ok {
		uid = fmt.Sprint(own.Uid)
		gid = fmt.Sprint(own.Gid)
		return uid, gid, owner
	}

	return "", "", ""

}

func md5checksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	h.Write([]byte(path))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash), nil
}
