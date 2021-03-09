package backup

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/wade-sam/fypclient/entity"
)

type Service struct {
	snrepo  SNRepository
	wffrepo WFRepository
	bsrepo  BSRepository
}

func NewBackupService(sn SNRepository, wf WFRepository, bs BSRepository) *Service {
	return &Service{
		snrepo:  sn,
		wffrepo: wf,
		bsrepo:  bs,
	}
}

func (s *Service) StartFullBackup(start, id string, ignorelist []string) error {
	ignoremap := make(map[string]string)
	for i := range ignorelist {
		ignoremap[ignorelist[i]] = ignorelist[i]
	}
	bsscan, err := s.DirectoryScan(start, ignoremap)
	if err != nil {
		return err
	}
	err = s.bsrepo.ExpectedFiles(bsscan)
	if err != nil {
		return err
	}

	return nil
}

func FullBackup() error {

	return nil
}

func (s *Service) StartIncrementalBackup(start, id string, ignorelist []string) error {
	ignoremap := make(map[string]string)
	for i := range ignorelist {
		ignoremap[ignorelist[i]] = ignorelist[i]
	}
	bsscan, err := s.DirectoryScan(start, ignoremap)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DirectoryScan(head string, ignore map[string]string) (*entity.Directory, error) {
	head = path.Clean(head)
	var bstree *entity.Directory
	var nodes = map[string]interface{}{}
	var walkfunc filepath.WalkFunc = func(loc string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
		} else if info.IsDir() {
			if _, ok := ignore[loc]; ok {
				return filepath.SkipDir
			} else {
				n := entity.NewDirectory(path.Base(loc))
				owner, uid, gid := GetProperties(info)
				n.AddProperties(owner, uid, gid)
				nodes[loc] = n
			}
			fmt.Println(loc)
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

			} else {
				return err
			}
		}

		return nil
	}
	err := filepath.Walk(head, walkfunc)
	if err != nil {
		return nil, err
	}

	for key, value := range nodes {
		var bsroot *entity.Directory
		if key == head {
			bstree = value.(*entity.Directory)
			continue
		} else {
			bsroot = nodes[path.Dir(key)].(*entity.Directory)
		}
		switch v := value.(type) {
		case *entity.Directory:
			bsroot.Folders[v.Name] = v
		case *entity.File:
			bsroot.Files = append(bsroot.Files, v)
		}
	}

	return bstree, nil
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
