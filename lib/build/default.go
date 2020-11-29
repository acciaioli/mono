package build

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/acciaioli/mono/internal/common"
	"github.com/acciaioli/mono/internal/specification"
	"github.com/acciaioli/mono/lib"
	"github.com/acciaioli/mono/lib/list"
)

const (
	BuildsRoot = ".builds"
)

type Service struct {
	Service      lib.Service
	ArtifactPath string
}

func BuildOutdated(bs common.BlobStorage) ([]Service, error) {
	services, err := list.List(bs)
	if err != nil {
		return nil, err
	}

	var diffServices []list.ListedService
	for _, service := range services {
		if service.Status == list.StatusDiff {
			diffServices = append(diffServices, service)
		}
	}
	return BuildArtifacts(diffServices)
}

func BuildArtifacts(services []list.ListedService) ([]Service, error) {
	var artifacts []Service
	for _, service := range services {
		artifactPath, err := BuildArtifact(&service)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, Service{Service: service.Service, ArtifactPath: *artifactPath})
	}
	return artifacts, nil
}

func BuildArtifact(service *list.ListedService) (*string, error) {
	artifactLocalPath, err := buildArtifact(
		service.Service.Path, &service.Service.Spec.Build.Artifact,
	)
	if err != nil {
		return nil, err
	}

	artifactBuildPath, err := moveArtifact(
		service.Service.Path, service.Checksum, *artifactLocalPath,
	)
	if err != nil {
		return nil, err
	}

	return artifactBuildPath, nil
}

func Clean() error {
	return os.RemoveAll(BuildsRoot)
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

func moveArtifact(service string, chsum string, localPath string) (*string, error) {
	buildPath := filepath.Join(BuildsRoot, service, fmt.Sprintf("%s.zip", chsum))
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
