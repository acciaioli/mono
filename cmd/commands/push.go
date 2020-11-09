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
			var pushed []push.Pushed
			var err error

			bucket, err := env.LoadArtifactBucket()
			if err != nil {
				return err
			}

			if artifact == "" {
				pushed, err = push.PushAllArtifacts(bucket)
				if err != nil {
					return err
				}
			} else {
				p, err := push.PushArtifact(bucket, artifact)
				if err != nil {
					return err
				}
				pushed = append(pushed, *p)
			}

			// todo: proper display
			for _, p := range pushed {
				fmt.Printf("Artifact: %s\n", p.Artifact)
				fmt.Printf("Status: %s\n", p.Status)
				if p.Err != nil {
					fmt.Printf("Error: %s\n", p.Err.Error())
				}
				if p.Key != nil {
					fmt.Printf("Key: %s\n", *p.Key)
				}
				fmt.Printf("\n")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&artifact, artifactFlag, "", artifactDescription)

	return cmd
}
