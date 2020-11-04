package aws

import (
	"github.com/pkg/errors"

	aws_session "github.com/aws/aws-sdk-go/aws/session"
)

func NewSession() (*aws_session.Session, error) {
	session, err := aws_session.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "error creating aws session")
	}

	return session, nil
}
