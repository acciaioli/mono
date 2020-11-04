package common

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type ServiceSpec struct {
	Build ServiceSpecBuild `yaml:"build"`
}

type ServiceSpecBuild struct {
	Artifact ServiceSpecBuildArtifact `yaml:"artifact"`
}

type ServiceSpecBuildArtifact struct {
	Name    string   `yaml:"name"`
	Command []string `yaml:"command"`
}

func LoadServiceSpec(servicePath string) (*ServiceSpec, error) {
	specFile := filepath.Join(servicePath, ServiceSpecFile)

	yamlB, err := ioutil.ReadFile(specFile)
	if err != nil {
		return nil, errors.Wrap(err, "error reading file")
	}

	spec := ServiceSpec{}
	err = yaml.Unmarshal(yamlB, &spec)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshal spec file")
	}

	return &spec, nil
}
