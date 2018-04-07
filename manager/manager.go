package manager

import (
	"errors"

	"github.com/brandoneprice31/freedrive/config"
	"github.com/brandoneprice31/freedrive/service"
)

type (
	Manager interface {
		Backup(string)
		Download(string, string)
	}

	manager struct {
		config   config.Config
		n        int
		services map[service.ServiceType]service.Service
	}

	services struct {
		services []service.Service
	}

	key struct {
		path              string
		storedServiceData []storedServiceData
	}

	storedServiceData struct {
		ServiceType service.ServiceType `json:"service_type"`
		ServiceData []byte              `json:"service_data"`
	}

	directory struct {
		directories []directory
		files       []file
		Path        string `json:"directory_path"`
	}

	file struct {
		Data []byte `json:"file_data"`
		size int
		Path string `json:"file_path"`
	}
)

var (
	ErrNoKey = errors.New("no key")
)

func NewManager(c config.Config, ss ...service.Service) Manager {
	ssMap := make(map[service.ServiceType]service.Service)
	for i := range ss {
		ssMap[ss[i].Type()] = ss[i]
	}
	return &manager{
		config:   c,
		n:        len(ss),
		services: ssMap,
	}
}
