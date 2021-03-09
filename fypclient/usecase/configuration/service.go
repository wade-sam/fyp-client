package configuration

import (
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"

	"github.com/wade-sam/fypclient/entity"
)

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
	head = path.Clean(head)
	var tree *entity.Directory
	var nodes = map[string]interface{}{}
	var walkfunc fs.WalkDirFunc = func(loc string, info fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
		} else if info.IsDir() {
			nodes[loc] = entity.NewDirectory(path.Base(loc))
		}
		return nil
	}
	err := filepath.WalkDir(head, walkfunc)
	if err != nil {
		return nil, entity.ErrFailedDirectoryScan
	}

	for key, value := range nodes {
		var root *entity.Directory
		if key == head {
			tree = value.(*entity.Directory)
			continue
		} else {
			root = nodes[path.Dir(key)].(*entity.Directory)
		}
		switch v := value.(type) {
		case *entity.Directory:
			root.Folders[v.Name] = v
		}
	}

	return tree, nil
}
