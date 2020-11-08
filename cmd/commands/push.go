package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/acciaioli/mono/cmd/env"
	"github.com/acciaioli/mono/services/push"
)

func Push() *cobra.Command {
	const (
		commandUse         = "push"
		commandDescription = "Pushes a service artifact to the cloud"

		artifactFlag        = "artifact"
		artifactDescription = "relative path to the artifact to be pushed"
	)

	var artifact string

	cmd := &cobra.Command{
		Use:   commandUse,
		Short: commandDescription,
		RunE: func(cmd2 *cobra.Command, args []string) error {
			bucket, err := env.LoadArtifactBucket()
			if err != nil {
				return err
			}

			location, err := push.Push(artifact, bucket)
			if err != nil {
				return err
			}

			fmt.Println(*location)
			return nil
		},
	}

	cmd.Flags().StringVar(&artifact, artifactFlag, "", artifactDescription)

	return cmd
}
