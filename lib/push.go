package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/acciaioli/mono/internal/common"
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

func PushAllArtifacts(bs common.BlobStorage) ([]Pushed, error) {
	var artifacts []string

	_, err := os.Stat(buildsRoot)
	if err != nil {
		return nil, errors.New("no artifacts found")
	}

	err = filepath.Walk(buildsRoot, func(path string, info os.FileInfo, err error) error {
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

	return PushArtifacts(bs, artifacts), nil
}

func PushArtifacts(bs common.BlobStorage, artifacts []string) []Pushed {
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
	if artifact == "" {
		return nil, errors.New("artifact cannot be empty")
	}
	if !strings.HasPrefix(artifact, buildsRoot) {
		return nil, errors.New(fmt.Sprintf("artifact should be under the builds root directory (%s)", buildsRoot))
	}
	if !common.FileExists(artifact) {
		return nil, errors.New(fmt.Sprintf("artifact `%s` does not exist", artifact))
	}

	tmpKey, err := filepath.Rel(buildsRoot, artifact)
	if err != nil {
		return nil, errors.Wrap(err, "error handling artifact path root")
	}

	servicePath := filepath.Dir(tmpKey)
	chsumExt := strings.Split(filepath.Base(tmpKey), ".")
	if len(chsumExt) != 2 {
		return nil, errors.Wrap(err, "error handling artifact path base")
	}

	version, err := nextVersion(bs, servicePath)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadFile(artifact)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error reading artifact file `%s`", artifact))
	}

	newKey := common.KeyFromArtifact(&common.Artifact{
		ServicePath: servicePath,
		Version:     *version,
		Checksum:    chsumExt[0],
		Extension:   chsumExt[1],
	})

	err = bs.Put(newKey, body)
	if err != nil {
		return nil, err
	}

	return &newKey, nil
}

func nextVersion(bs common.BlobStorage, service string) (*int, error) {
	latestKey, err := bs.GetLatestKey(service)
	if err != nil {
		return nil, err
	}

	if latestKey == nil {
		one := 1
		return &one, nil
	} else {
		a, err := common.ArtifactFromKey(*latestKey)
		if err != nil {
			return nil, err
		}
		version := a.Version + 1
		return &version, nil
	}
}
