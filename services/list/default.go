package list

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/acciaioli/mono/services/checksum"

	"github.com/pkg/errors"

	"github.com/acciaioli/mono/internal/common"
)

type Service struct {
	Path                 string
	Status               *string
	Checksum             *string
	LatestPushedChecksum *string
}

func List(bucket string) ([]Service, error) {
	paths, err := listServices("./")
	if err != nil {
		return nil, err
	}

	return toServices(paths, bucket)
}

func listServices(root string) ([]string, error) {
	var services []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "error walking directory")
		}

		for _, service := range services {
			if strings.Contains(path, service) {
				return nil // we don't support nested services
			}
		}

		if common.IsServiceDir(path) {
			services = append(services, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return services, nil
}

func toServices(servicePaths []string, bucket string) ([]Service, error) {
	type serviceErr struct {
		s *Service
		e error
	}
	serviceErrChan := make(chan serviceErr, len(servicePaths))

	wg := sync.WaitGroup{}

	for _, path := range servicePaths {
		path := path
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := Service{Path: path}
			chsum, err := checksum.ComputeChecksum(path)
			if err != nil {
				serviceErrChan <- serviceErr{e: err}
				return
			}
			s.Checksum = chsum

			pushedChsum, err := checksum.GetLatestChecksum(path, bucket)
			if err != nil {
				serviceErrChan <- serviceErr{e: err}
				return
			}
			s.LatestPushedChecksum = pushedChsum

			status := func() string {
				if *chsum == *pushedChsum {
					return "ok"
				}
				return "diff"
			}()
			s.Status = &status

			serviceErrChan <- serviceErr{s: &s}
		}()
	}
	wg.Wait()
	close(serviceErrChan)

	var services []Service
	var allErrs []string
	for se := range serviceErrChan {
		if se.e != nil {
			allErrs = append(allErrs, se.e.Error())
			continue
		}
		services = append(services, *se.s)
	}

	if allErrs != nil {
		return nil, errors.New(fmt.Sprintf("errors: %s", strings.Join(allErrs, ";")))
	}

	return services, nil
}
