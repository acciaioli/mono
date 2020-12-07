package commands

import (
	"github.com/spf13/cobra"

	"github.com/acciaioli/mono/cmd/display"
	"github.com/acciaioli/mono/cmd/env"
	"github.com/acciaioli/mono/internal/common"
	"github.com/acciaioli/mono/lib"
)

func Push() *cobra.Command {
	const (
		commandUse         = "push"
		commandDescription = "Pushes service artifacts to the cloud"

		artifactsFlag        = "artifact"
		artifactsDescription = "relative path to the artifact(s) to be pushed"

		keepArtifactsFlag        = "keep"
		keepArtifactsDescription = "keep successfully uploaded artifacts"
	)

	var artifacts []string
	var keep bool

	cmd := &cobra.Command{
		Use:   commandUse,
		Short: commandDescription,
		RunE: func(cmd2 *cobra.Command, args []string) error {
			bucket, err := env.LoadArtifactBucket()
			if err != nil {
				return err
			}
			bs, err := common.NewBlobStorage(bucket)
			if err != nil {
				return err
			}

			pushed, err := func() ([]lib.Pushed, error) {
				if artifacts == nil {
					return lib.PushAllArtifacts(bs, keep)
				}
				return lib.PushArtifacts(bs, artifacts, keep), nil
			}()
			if err != nil {
				return err
			}

			headers := []string{"artifact", "status", "key", "error"}
			var data [][]string
			for _, row := range pushed {
				if row.Err != nil {
					data = append(data, []string{row.Artifact, string(row.Status), "-", row.Err.Error()})
				}
				data = append(data, []string{row.Artifact, string(row.Status), *row.Key, "-"})
			}
			display.Table(headers, data)

			return nil
		},
	}

	cmd.Flags().StringArrayVar(&artifacts, artifactsFlag, nil, artifactsDescription)
	cmd.Flags().BoolVar(&keep, keepArtifactsFlag, false, keepArtifactsDescription)

	return cmd
}
