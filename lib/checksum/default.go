package checksum

import (
	"crypto/sha1" // nolint
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/acciaioli/mono/internal/common"
	"github.com/acciaioli/mono/internal/patterns"
	"github.com/acciaioli/mono/lib"
)

func ComputeChecksum(service *lib.Service) (*string, error) {
	return computeChecksum(service.Path, service.Spec.Checksum.Exclude)
}

func GetLatestChecksum(service *lib.Service, bs common.BlobStorage) (*string, error) {
	return getLatestChecksum(service.Path, bs)
}

func computeChecksum(servicePath string, excludedPatterns []string) (*string, error) {
	hash := sha1.New() // nolint

	excluded, err := patterns.NewMatcher(excludedPatterns)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(servicePath, func(fPath string, fInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(servicePath, fPath)
		if err != nil {
			return err
		}

		if excluded.Match(rel) {
			return nil
		}

		if fInfo.IsDir() {
			return nil
		}

		fReader, err := os.Open(fPath)
		if err != nil {
			return err
		}

		_, err = hash.Write([]byte(fPath))
		if err != nil {
			return err
		}

		content, err := ioutil.ReadAll(fReader)
		if err != nil {
			return err
		}

		_, err = hash.Write(content)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s := fmt.Sprintf("%x", hash.Sum(nil))
	return &s, nil
}

func getLatestChecksum(servicePath string, bs common.BlobStorage) (*string, error) {
	latestKey, err := bs.GetLatestKey(servicePath)
	if err != nil {
		return nil, err
	}
	if latestKey == nil {
		chsum := ""
		return &chsum, nil
	}
	artifact, err := common.SplitArtifactKey(*latestKey)
	if err != nil {
		return nil, err
	}

	return &artifact.Checksum, nil
}
