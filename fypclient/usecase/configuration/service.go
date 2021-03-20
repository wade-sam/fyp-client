package configuration

import (
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"

	jsoniter "github.com/json-iterator/go"

	"github.com/wade-sam/fypclient/entity"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Config struct {
	ClientName string   `json:"clientname"`
	Policies   []string `json:"policies"`
}
type Service struct {
	repo Repository
}

func NewConfigurationService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) GetClientName() (string, error) {
	//for hostname check existing host isn't empty. If it is then get the hostname of the system and append random characters on the end e.g. the time/date
	name, err := s.repo.GetClientName()
	if err != nil {
		return "", err
	}
	if name == "" {
		id := rand.Intn(5)
		prename, err := os.Hostname()
		newname := fmt.Sprintf("%v%v", prename, id)
		err = s.repo.SetClientName(newname)
		if err != nil {
			return "", err
		}
		return newname, nil
	}
	return name, nil
}

func (s *Service) SetStorageNode(name string) error {
	return s.repo.SetStorageNode(name)
}

func (s *Service) GetStorageNode() (string, error) {
	name, err := s.repo.GetStorageNode()
	if err != nil {
		return "", err
	}
	if name == "" {
		return "", entity.ErrFieldWasEmpty
	}
	return name, nil
}

func (s *Service) DirectoryScan(head string) (*entity.Directory, error) {
	head = path.Clean(head) //cleans up the path, returning the shortest filepath
	var tree *entity.Directory
	var nodes = map[string]interface{}{}
	var walkfunc fs.WalkDirFunc = func(loc string, info fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
		} else if info.IsDir() {
			nodes[loc] = entity.NewDirectory(path.Base(loc)) //loc=filepath path.Base(loc) = dir name
		}
		return nil
	}
	err := filepath.WalkDir(head, walkfunc)
	if err != nil {
		return nil, entity.ErrFailedDirectoryScan
	}
	// this keeps looping through until all the directories have been placed inside eachother.
	for key, value := range nodes {
		var root *entity.Directory // initialise root as a pointer of entity.Directory. Doesn't have an actual value yet
		if key == head {
			tree = value.(*entity.Directory) //if key is equal to head then assign value to tree
			fmt.Println(tree)
			continue
		} else {
			root = nodes[path.Dir(key)].(*entity.Directory) //root is NOW assigned as the parent directory of values directory. path.Dir takes off the last element in the path
		}
		switch v := value.(type) {
		case *entity.Directory:
			root.Folders[v.Name] = v //add value to the parents folders.
		}
	}

	return tree, nil
}

func (s *Service) WriteBackupResult(files map[string]*entity.FileDTO) error {
	err := s.repo.CreateBackupResult(files)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetBackupResult() (map[string]*entity.FileDTO, error) {
	return nil, nil
}
