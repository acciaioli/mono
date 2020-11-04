package aws

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	aws_s3 "github.com/aws/aws-sdk-go/service/s3"
	aws_s3_manager "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3BlobStorage struct {
	session *aws_session.Session
	bucket  string
}

func NewS3BlobStorage(bucket string) (*S3BlobStorage, error) {
	session, err := NewSession()
	if err != nil {
		return nil, err
	}
	return &S3BlobStorage{session: session, bucket: bucket}, nil
}

func (bs *S3BlobStorage) Get(key string) ([]byte, error) {
	return S3Download(bs.session, bs.bucket, key)
}

func (bs *S3BlobStorage) Put(key string, data []byte) error {
	return S3Upload(bs.session, bs.bucket, key, data)
}

func (bs *S3BlobStorage) GetLatestKey(prefix string) (*string, error) {
	return S3GetLatestKey(bs.session, bs.bucket, prefix)
}

func S3Download(session *aws_session.Session, bucket, key string) ([]byte, error) {
	downloader := aws_s3_manager.NewDownloader(session)

	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(buf, &aws_s3.GetObjectInput{
		Bucket: aws.String(bucket), Key: aws.String(key),
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error downloading `%s` from aws s3", filepath.Join(bucket, key)))
	}

	return buf.Bytes(), nil
}

func S3Upload(session *aws_session.Session, bucket, key string, data []byte) error {
	uploader := aws_s3_manager.NewUploader(session)

	_, err := uploader.Upload(&aws_s3_manager.UploadInput{
		Bucket: aws.String(bucket), Key: aws.String(key), Body: bytes.NewReader(data),
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error uploading `%s` to aws s3", filepath.Join(bucket, key)))
	}

	return nil
}

func S3GetLatestKey(session *aws_session.Session, bucket string, prefix string) (*string, error) {
	s3 := aws_s3.New(session)

	listResp, err := s3.ListObjectsV2(&aws_s3.ListObjectsV2Input{
		Bucket: aws.String(bucket), Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error listing aws s3 objects from bucket `%s`", bucket))
	}
	if len(listResp.Contents) < 1 {
		return nil, nil
	}

	latest := listResp.Contents[0]

	for _, obj := range listResp.Contents {
		if obj.LastModified.After(*latest.LastModified) {
			latest = obj
		}
	}

	return latest.Key, nil
}
