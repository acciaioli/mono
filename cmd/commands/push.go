package commands

import (
	"github.com/spf13/cobra"

	"github.com/acciaioli/mono/cmd/display"
	"github.com/acciaioli/mono/cmd/env"
	"github.com/acciaioli/mono/internal/common"
	"github.com/acciaioli/mono/lib/push"
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
			bs, err := common.NewBlobStorage(bucket)
			if err != nil {
				return err
			}

			if artifact == "" {
				pushed, err = push.PushAllArtifacts(bs)
				if err != nil {
					return err
				}
			} else {
				p, err := push.PushArtifact(bs, artifact)
				if err != nil {
					return err
				}
				pushed = append(pushed, *p)
			}

			headers := []string{"artifact", "status", "key", "error"}
			var data [][]string
			for _, p := range pushed {
				if p.Err != nil {
					data = append(data, []string{p.Artifact, string(p.Status), "", p.Err.Error()})
				}
				data = append(data, []string{p.Artifact, string(p.Status), *p.Key, ""})
			}
			display.Table(headers, data)

			return nil
		},
	}

	cmd.Flags().StringVar(&artifact, artifactFlag, "", artifactDescription)

	return cmd
}
