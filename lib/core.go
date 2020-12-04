package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/acciaioli/mono/internal/common"
	"github.com/acciaioli/mono/internal/specification"
)

type Service struct {
	Path string
	Spec specification.Spec
}

type loadSettings struct {
	root          string
	specificPaths []string
	// filters
	containing   *string
	startingWith *string
}

func defaultSettings() *loadSettings {
	return &loadSettings{
		root: "./",
	}
}

type loadOption func(s *loadSettings)

func loadSpecificServices(paths []string) loadOption {
	return func(s *loadSettings) {
		s.specificPaths = paths
	}
}

func LoadService(path string) (*Service, error) {
	if path == "" {
		return nil, errors.New("service path cannot be empty")
	}
	if !common.IsServiceDir(path) {
		return nil, errors.New(fmt.Sprintf("%s is not a valid service path", path))
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

func LoadServices(opts ...loadOption) ([]Service, error) {
	s := defaultSettings()
	for _, opt := range opts {
		opt(s)
	}

	paths, err := func() ([]string, error) {
		if s.specificPaths != nil {
			return s.specificPaths, nil
		} else {
			return scan(s.root)
		}
	}()
	if err != nil {
		return nil, err
	}
	paths = filter(s, paths)
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

func scan(root string) ([]string, error) {
	var servicePaths []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "error walking directory")
		}

		for _, service := range servicePaths {
			if strings.Contains(path, service) {
				return nil // we don't support nested servicePaths
			}
		}

		if common.IsServiceDir(path) {
			servicePaths = append(servicePaths, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return servicePaths, nil
}

func filter(s *loadSettings, paths []string) []string {
	if s.containing != nil {
		var fpaths []string
		for _, path := range paths {
			if strings.Contains(path, *s.containing) {
				fpaths = append(fpaths, path)
			}
		}
		paths = fpaths
	}

	if s.startingWith != nil {
		var fpaths []string
		for _, path := range paths {
			if strings.HasPrefix(path, *s.startingWith) {
				fpaths = append(fpaths, path)
			}
		}
		paths = fpaths
	}

	return paths
}
