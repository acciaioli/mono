package lib

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/acciaioli/mono/internal/common"
	"github.com/acciaioli/mono/internal/specification"
)

const (
	buildsRoot = ".builds"
)

type Build struct {
	Service      Service
	ArtifactPath string
}

func BuildOutdatedServices(bs common.BlobStorage) ([]Build, error) {
	services, err := LoadServices()
	if err != nil {
		return nil, err
	}

	states, err := GetServicesState(bs, services)
	if err != nil {
		return nil, err
	}

	var diffs []ServiceState
	for _, state := range states {
		if state.Diff() {
			diffs = append(diffs, state)
		}
	}
	return buildServices(diffs)
}

func BuildServices(bs common.BlobStorage, servicePaths []string) ([]Build, error) {
	services, err := LoadServices(loadSpecificServices(servicePaths))
	if err != nil {
		return nil, err
	}

	states, err := GetServicesState(bs, services)
	if err != nil {
		return nil, err
	}

	return buildServices(states)
}

func Clean() error {
	return os.RemoveAll(buildsRoot)
}

func buildServices(states []ServiceState) ([]Build, error) {
	var builds []Build
	for _, lService := range states {
		artifactPath, err := buildService(&lService)
		if err != nil {
			return nil, err
		}
		builds = append(builds, Build{Service: *lService.Service, ArtifactPath: *artifactPath})
	}
	return builds, nil
}

func buildService(state *ServiceState) (*string, error) {
	artifactLocalPath, err := buildArtifact(
		state.Service.Path, &state.Service.Spec.Build.Artifact,
	)
	if err != nil {
		return nil, err
	}

	artifactBuildPath, err := moveArtifact(
		state.Service.Path, state.LocalChecksum, *artifactLocalPath,
	)
	if err != nil {
		return nil, err
	}

	return artifactBuildPath, nil
}

func buildArtifact(servicePath string, artifactSpec *specification.BuildArtifact) (*string, error) {
	if len(artifactSpec.Command) < 1 {
		return nil, errors.New("build artifact has length 0")
	}
	cmd := exec.Command(artifactSpec.Command[0], artifactSpec.Command[1:]...) // nolint
	//log.Printf("running command: %s\n", cmd.String())
	cmd.Dir = servicePath
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	//log.Println("command output:")
	//log.Println(stdout.String())
	if err != nil {
		return nil, errors.Wrap(err, "error running build artifact command")
	}

	artifactPath := filepath.Join(servicePath, artifactSpec.Name)
	if !common.FileExists(artifactPath) {
		return nil, errors.New("build artifact not found after successful build artifact command run")
	}

	return &artifactPath, nil
}

func moveArtifact(servicePath string, chsum string, localPath string) (*string, error) {
	buildPath := filepath.Join(buildsRoot, servicePath, fmt.Sprintf("%s.zip", chsum))

	err := os.MkdirAll(filepath.Dir(buildPath), os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "error renaming file")
	}

	err = os.Rename(localPath, buildPath)
	if err != nil {
		return nil, errors.Wrap(err, "error renaming file")
	}

	return &buildPath, nil
}
