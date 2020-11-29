package list

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/acciaioli/mono/internal/common"
	"github.com/acciaioli/mono/lib"
	"github.com/acciaioli/mono/lib/checksum"
)

type Status string

const (
	StatusOK   Status = "ok"
	StatusDiff Status = "diff"
)

type ListedService struct {
	Service              lib.Service
	Status               Status
	Checksum             string
	LatestPushedChecksum string
}

func List(bs common.BlobStorage) ([]ListedService, error) {
	paths, err := discoverServicePaths("./")
	if err != nil {
		return nil, err
	}

	services, err := lib.LoadServices(paths)
	if err != nil {
		return nil, err
	}
	return servicesStatus(services, bs)
}

func discoverServicePaths(root string) ([]string, error) {
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

func servicesStatus(services []lib.Service, bs common.BlobStorage) ([]ListedService, error) {
	type serviceErr struct {
		s *ListedService
		e error
	}
	serviceErrChan := make(chan serviceErr, len(services))

	wg := sync.WaitGroup{}

	for _, service := range services {
		service := service
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := ListedService{Service: service}
			chsum, err := checksum.ComputeChecksum(&service)
			if err != nil {
				serviceErrChan <- serviceErr{e: err}
				return
			}
			s.Checksum = *chsum

			pushedChsum, err := checksum.GetLatestChecksum(&service, bs)
			if err != nil {
				serviceErrChan <- serviceErr{e: err}
				return
			}
			s.LatestPushedChecksum = *pushedChsum

			s.Status = func() Status {
				if *chsum == *pushedChsum {
					return StatusOK
				}
				return StatusDiff
			}()

			serviceErrChan <- serviceErr{s: &s}
		}()
	}
	wg.Wait()
	close(serviceErrChan)

	var listedServices []ListedService
	var allErrs []string
	for se := range serviceErrChan {
		if se.e != nil {
			allErrs = append(allErrs, se.e.Error())
			continue
		}
		listedServices = append(listedServices, *se.s)
	}

	if allErrs != nil {
		return nil, errors.New(fmt.Sprintf("errors:\n%s", strings.Join(allErrs, "\n")))
	}

	return listedServices, nil
}
