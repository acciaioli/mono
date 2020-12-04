package common

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Artifact struct {
	ServicePath string
	Version     int
	Checksum    string
	Extension   string
}

func ArtifactFromKey(key string) (*Artifact, error) {
	servicePath := filepath.Dir(key)
	splitBase := strings.Split(filepath.Base(key), ".")
	if len(splitBase) != 3 {
		return nil, errors.New("error parsing artifact key")
	}

	version, err := strconv.Atoi(splitBase[0][1:])
	if err != nil {
		return nil, errors.Wrap(err, "error converting artifact version to int")
	}

	return &Artifact{
		ServicePath: servicePath,
		Version:     version,
		Checksum:    splitBase[1],
		Extension:   splitBase[2],
	}, nil
}

func KeyFromArtifact(artifact *Artifact) string {
	return filepath.Join(
		artifact.ServicePath,
		fmt.Sprintf("v%d.%s.%s", artifact.Version, artifact.Checksum, artifact.Extension),
	)
}
