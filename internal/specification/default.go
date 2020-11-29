package specification

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Spec struct {
	Checksum Checksum `yaml:"checksum"`
	Build    Build    `yaml:"build"`
}

type Checksum struct {
	Exclude []string `yaml:"exclude"`
}

type Build struct {
	Artifact BuildArtifact `yaml:"artifact"`
}

type BuildArtifact struct {
	Name    string   `yaml:"name"`
	Command []string `yaml:"command"`
}

func Load(specPath string) (*Spec, error) {
	yamlB, err := ioutil.ReadFile(specPath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading spec file")
	}

	spec := Spec{}
	err = yaml.Unmarshal(yamlB, &spec)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling spec file")
	}

	return &spec, nil
}
