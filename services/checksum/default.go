package checksum

import (
	"crypto/sha1" // nolint
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/acciaioli/mono/internal/specification"

	"github.com/acciaioli/mono/internal/patterns"

	"github.com/pkg/errors"

	"github.com/acciaioli/mono/internal/common"
)

func ComputeChecksum(service string) (*string, error) {
	if service == "" {
		return nil, errors.New("service cannot be empty")
	}
	if !common.IsServiceDir(service) {
		return nil, errors.New(fmt.Sprintf("%s is not a valid service path", service))
	}
	serviceSpec, err := specification.Load(service)
	if err != nil {
		return nil, err
	}

	return computeChecksum(service, serviceSpec.Checksum.Exclude)
}

func GetLatestChecksum(service string, bucket string) (*string, error) {
	if service == "" {
		return nil, errors.New("service cannot be empty")
	}
	if !common.IsServiceDir(service) {
		return nil, errors.New(fmt.Sprintf("%s is not a valid service path", service))
	}
	if bucket == "" {
		return nil, errors.New("bucket cannot be empty")
	}

	bs, err := common.NewBlobStorage(bucket)
	if err != nil {
		return nil, err
	}

	return getLatestChecksum(service, bs)
}

func computeChecksum(serviceRoot string, excludedPatterns []string) (*string, error) {
	hash := sha1.New() // nolint

	excluded, err := patterns.NewMatcher(excludedPatterns)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(serviceRoot, func(fPath string, fInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(serviceRoot, fPath)
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

func getLatestChecksum(service string, bs common.BlobStorage) (*string, error) {
	latestKey, err := bs.GetLatestKey(service)
	if err != nil {
		return nil, err
	}
	if latestKey == nil {
		chsum := ""
		return &chsum, nil
	}
	chsum := strings.TrimSuffix(filepath.Base(*latestKey), ".zip")
	return &chsum, nil
}
