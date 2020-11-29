package lib

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/acciaioli/mono/internal/common"
	"github.com/acciaioli/mono/internal/specification"
)

type Service struct {
	Path string
	Spec specification.Spec
}

func LoadService(path string) (*Service, error) {
	if path == "" {
		return nil, errors.New("service path cannot be empty")
	}
	if !common.IsServiceDir(path) {
		return nil, errors.New(fmt.Sprintf("%s is not a service path", path))
	}

	specPath := filepath.Join(path, common.ServiceSpecFile)
	spec, err := specification.Load(specPath)
	if err != nil {
		return nil, err
	}

	return &Service{
		Path: path,
		Spec: *spec,
	}, nil
}

func LoadServices(paths []string) ([]Service, error) {
	var services []Service
	for _, path := range paths {
		service, err := LoadService(path)
		if err != nil {
			return nil, err
		}
		services = append(services, *service)
	}
	return services, nil
}
