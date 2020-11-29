package common

import (
	"errors"
	"path/filepath"
	"strings"
)

type Artifact struct {
	Version  string
	Checksum string
}

func BuildArtifactKey() {

}

func SplitArtifactKey(key string) (*Artifact, error) {
	versionChsum := strings.Split(strings.TrimSuffix(filepath.Base(key), ".zip"), ".")
	if len(versionChsum) != 2 {
		return nil, errors.New("error parsing artifact key")
	}
	return &Artifact{
		Version:  versionChsum[0],
		Checksum: versionChsum[1],
	}, nil
}
