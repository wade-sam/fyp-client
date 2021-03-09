package configuration

import (
	"github.com/wade-sam/fypclient/entity"
)

type Repository interface {
	GetClientName() (string, error)
	SetClientName(name string) error
	SetStorageNode(name string) error
	GetStorageNode() (string, error)
}

type Usecase interface {
	GetClientName() (string, error)
	GetStorageNode() (string, error)
	SetStorageNode(name string) error
	DirectoryScan(start string) (*entity.Directory, error)
}
