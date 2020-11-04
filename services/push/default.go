package push

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"

	"github.com/acciaioli/mono/internal/common"
)

func Push(artifact string, bucket string) (*string, error) {
	err := validate(artifact, bucket)
	if err != nil {
		return nil, err
	}

	bs, err := common.NewBlobStorage(bucket)
	if err != nil {
		return nil, err
	}

	return push(artifact, bs)
}

func validate(artifact string, bucket string) error {
	if artifact == "" {
		return errors.New("artifact cannot be empty")
	}
	if !strings.HasPrefix(artifact, common.BuildsRoot) {
		return errors.New(fmt.Sprintf("artifact should be under the builds root directory (%s)", common.BuildsRoot))
	}
	if !common.FileExists(artifact) {
		return errors.New(fmt.Sprintf("artifact `%s` does not exist", artifact))
	}

	if bucket == "" {
		return errors.New("bucket cannot be empty")
	}
	return nil
}

func push(artifact string, bs common.BlobStorage) (*string, error) {
	body, err := ioutil.ReadFile(artifact)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error reading artifact file `%s`", artifact))
	}

	key := strings.TrimPrefix(artifact, common.BuildsRoot)
	err = bs.Put(key, body)
	if err != nil {
		return nil, err
	}

	return &key, nil
}
