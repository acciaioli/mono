package build

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/acciaioli/mono/services/checksum"

	"github.com/pkg/errors"

	"github.com/acciaioli/mono/internal/common"
)

func Build(service string) (*string, error) {
	if !common.IsServiceDir(service) {
		return nil, errors.New(fmt.Sprintf("%s is not a valid service path", service))
	}

	spec, err := common.LoadServiceSpec(service)
	if err != nil {
		return nil, err
	}

	chsum, err := checksum.ComputeChecksum(service)
	if err != nil {
		return nil, err
	}

	artifactLocalPath, err := buildArtifact(service, &spec.Build.Artifact)
	if err != nil {
		return nil, err
	}

	artifactBuildPath, err := moveArtifact(service, *chsum, *artifactLocalPath)
	if err != nil {
		return nil, err
	}

	return artifactBuildPath, nil
}

func buildArtifact(serviceDir string, artifactSpec *common.ServiceSpecBuildArtifact) (*string, error) {
	if len(artifactSpec.Command) < 1 {
		return nil, errors.New("build artifact has length 0")
	}
	cmd := exec.Command(artifactSpec.Command[0], artifactSpec.Command[1:]...)
	//log.Printf("running command: %s\n", cmd.String())
	cmd.Dir = serviceDir
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

	artifactPath := filepath.Join(serviceDir, artifactSpec.Name)
	if !common.FileExists(artifactPath) {
		return nil, errors.New("build artifact not found after successful build artifact command run")
	}

	return &artifactPath, nil
}

func moveArtifact(service string, chsum string, localPath string) (*string, error) {
	buildPath := filepath.Join(common.BuildsRoot, service, fmt.Sprintf("%s.zip", chsum))
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
