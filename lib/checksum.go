package lib

import (
	"crypto/sha1" // nolint
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/acciaioli/mono/internal/patterns"
)

type Checksum struct {
	Service  *Service
	Checksum string
}

func ComputeChecksums(servicePaths []string) ([]Checksum, error) {
	services, err := LoadServices(loadSpecificServices(servicePaths))
	if err != nil {
		return nil, err
	}

	var checksums []Checksum
	for _, service := range services {
		service := service
		checksum, err := computeChecksum(service.Path, service.Spec.Checksum.Exclude)
		if err != nil {
			return nil, err
		}
		checksums = append(checksums, Checksum{
			Service:  &service,
			Checksum: checksum,
		})
	}
	return checksums, nil
}

func computeChecksum(servicePath string, excludedPatterns []string) (string, error) {
	hash := sha1.New() // nolint

	excluded, err := patterns.NewMatcher(excludedPatterns)
	if err != nil {
		return "", err
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
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
