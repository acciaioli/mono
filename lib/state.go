package lib

import (
	"errors"
	"strings"
	"sync"

	"github.com/acciaioli/mono/internal/common"
)

type ServiceState struct {
	Service       *Service
	LocalChecksum string
	Checksum      string
	Version       int
}

func (s ServiceState) Diff() bool {
	return s.LocalChecksum != s.Checksum
}

func GetServiceState(bs common.BlobStorage, service *Service) (*ServiceState, error) {
	localChecksum, err := computeChecksum(service.Path, service.Spec.Checksum.Exclude)
	if err != nil {
		return nil, err
	}

	artifact, err := getLatestArtifact(bs, service)
	if err != nil {
		return nil, err
	}

	checksum, version := func() (string, int) {
		if artifact == nil {
			return "-", 0
		}
		return artifact.Checksum, artifact.Version
	}()

	return &ServiceState{
		Service:       service,
		LocalChecksum: localChecksum,
		Checksum:      checksum,
		Version:       version,
	}, nil
}

func GetServicesState(bs common.BlobStorage, services []Service) ([]ServiceState, error) {
	type stateOrError struct {
		serviceState *ServiceState
		err          error
	}
	channel := make(chan stateOrError, len(services))

	wg := sync.WaitGroup{}

	for _, service := range services {
		service := service
		wg.Add(1)
		go func() {
			defer wg.Done()
			serviceState, err := GetServiceState(bs, &service)
			channel <- stateOrError{serviceState: serviceState, err: err}
		}()
	}
	wg.Wait()
	close(channel)

	var states []ServiceState
	var errs []error
	for msg := range channel {
		if msg.err != nil {
			errs = append(errs, msg.err)
		}
		states = append(states, *msg.serviceState)
	}
	if errs != nil {
		var errStrings []string
		for _, err := range errs {
			errStrings = append(errStrings, err.Error())
		}
		return nil, errors.New(strings.Join(errStrings, "\n"))
	}
	return states, nil
}

func getLatestArtifact(bs common.BlobStorage, service *Service) (*common.Artifact, error) {
	latestKey, err := bs.GetLatestKey(service.Path)
	if err != nil {
		return nil, err
	}
	if latestKey == nil {
		return nil, nil
	}
	return common.ArtifactFromKey(*latestKey)
}
