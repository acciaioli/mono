package push

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/acciaioli/mono/internal/common"
	"github.com/acciaioli/mono/lib/build"
)

type PushStatus string

const (
	StatusSuccessful PushStatus = "successful"
	StatusFailed     PushStatus = "failed"
)

type Pushed struct {
	Artifact string
	Status   PushStatus
	Key      *string
	Err      error
}

func PushAllArtifacts(bucket string) ([]Pushed, error) {
	var artifacts []string

	_, err := os.Stat(build.BuildsRoot)
	if err != nil {
		return nil, nil
	}

	err = filepath.Walk(build.BuildsRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "error walking builds directory")
		}

		if !info.IsDir() {
			artifacts = append(artifacts, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	bs, err := common.NewBlobStorage(bucket)
	if err != nil {
		return nil, err
	}

	return pushArtifacts(bs, artifacts), nil
}

func PushArtifact(bucket string, artifact string) (*Pushed, error) {
	if artifact == "" {
		return nil, errors.New("artifact cannot be empty")
	}
	if !strings.HasPrefix(artifact, build.BuildsRoot) {
		return nil, errors.New(fmt.Sprintf("artifact should be under the builds root directory (%s)", build.BuildsRoot))
	}
	if !common.FileExists(artifact) {
		return nil, errors.New(fmt.Sprintf("artifact `%s` does not exist", artifact))
	}

	bs, err := common.NewBlobStorage(bucket)
	if err != nil {
		return nil, err
	}

	key, err := pushArtifact(bs, artifact)
	if err != nil {
		return &Pushed{Artifact: artifact, Status: StatusFailed, Err: err}, nil
	}
	return &Pushed{Artifact: artifact, Status: StatusSuccessful, Key: key}, nil
}

func pushArtifacts(bs common.BlobStorage, artifacts []string) []Pushed {
	pushedErrChan := make(chan Pushed, len(artifacts))
	wg := sync.WaitGroup{}

	for _, artifact := range artifacts {
		artifact := artifact
		wg.Add(1)
		go func() {
			defer wg.Done()
			key, err := pushArtifact(bs, artifact)
			if err != nil {
				pushedErrChan <- Pushed{Artifact: artifact, Status: StatusFailed, Err: err}
				return
			}
			pushedErrChan <- Pushed{Artifact: artifact, Status: StatusSuccessful, Key: key}
		}()
	}
	wg.Wait()
	close(pushedErrChan)

	var pushed []Pushed
	for item := range pushedErrChan {
		pushed = append(pushed, item)
	}

	return pushed
}

func pushArtifact(bs common.BlobStorage, artifact string) (*string, error) {
	body, err := ioutil.ReadFile(artifact)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error reading artifact file `%s`", artifact))
	}

	key := strings.TrimPrefix(artifact, build.BuildsRoot)
	err = bs.Put(key, body)
	if err != nil {
		return nil, err
	}

	return &key, nil
}
