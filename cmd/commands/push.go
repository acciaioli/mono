package commands

import (
	"fmt"

	"github.com/acciaioli/mono/services/push"
	"github.com/spf13/cobra"
)

func Push() *cobra.Command {
	const (
		commandUse         = "push"
		commandDescription = "Pushes a service artifact to the cloud"

		artifactFlag        = "artifact"
		artifactDescription = "relative path to the artifact to be pushed"

		bucketFlag        = "bucket"
		bucketDescription = "artifacts bucket. format should be one of ['s3://',]"
	)

	var artifact string
	var bucket string

	cmd := &cobra.Command{
		Use:   commandUse,
		Short: commandDescription,
		RunE: func(cmd2 *cobra.Command, args []string) error {
			location, err := push.Push(artifact, bucket)
			if err != nil {
				return err
			}

			fmt.Println(*location)
			return nil
		},
	}

	cmd.Flags().StringVar(&artifact, artifactFlag, "", artifactDescription)
	cmd.Flags().StringVar(&bucket, bucketFlag, "", bucketDescription)

	return cmd
}
