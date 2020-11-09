package env

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

const (
	envArtifactBucket = "MONO_ARTIFACT_BUCKET"
)

func LoadArtifactBucket() (string, error) {
	bucket := os.Getenv(envArtifactBucket)
	if bucket == "" {
		return "", errors.New(fmt.Sprintf("env variable %s must be set", envArtifactBucket))
	}
	return bucket, nil
}
