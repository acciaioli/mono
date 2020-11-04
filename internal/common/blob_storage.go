package common

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/acciaioli/mono/internal/aws"
)

const (
	BucketPrefixS3 = "s3://"
)

type BlobStorage interface {
	Get(key string) ([]byte, error)
	Put(key string, data []byte) error
	GetLatestKey(prefix string) (*string, error)
}

func NewBlobStorage(bucket string) (BlobStorage, error) {
	prefix, ok := isBucketSupported(bucket)
	if !ok {
		return nil, errors.New(fmt.Sprintf("bucket `%s` is not supported", bucket))
	}

	var bs BlobStorage
	var err error

	switch prefix {
	case BucketPrefixS3:
		bs, err = aws.NewS3BlobStorage(strings.TrimPrefix(bucket, prefix))
	default:
		return nil, errors.New(fmt.Sprintf("bucket prefix `%s` doesn't have registered implementation", prefix))
	}

	if err != nil {
		return nil, errors.Wrap(err, "error creating blob storage")
	}

	return bs, nil
}

func isBucketSupported(bucket string) (string, bool) {
	for _, prefix := range []string{BucketPrefixS3} {
		if strings.HasPrefix(bucket, prefix) {
			return prefix, true
		}
	}
	return "", false
}
