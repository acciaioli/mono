package commands

import (
	"github.com/spf13/cobra"

	"github.com/acciaioli/mono/cmd/display"
	"github.com/acciaioli/mono/cmd/env"
	"github.com/acciaioli/mono/internal/common"
	"github.com/acciaioli/mono/lib"
)

func Build() *cobra.Command {
	const (
		commandUse         = "build"
		commandDescription = "Builds artifact for a service"

		servicesFlag        = "service"
		servicesDescription = "relative path to the target service root directory"

		cleanFlag        = "clean"
		cleanDescription = "cleans all built artifacts"
	)

	var servicePaths []string
	var clean bool

	cmd := &cobra.Command{
		Use:   commandUse,
		Short: commandDescription,
		RunE: func(cmd2 *cobra.Command, args []string) error {
			if clean {
				return lib.Clean()
			}

			bucket, err := env.LoadArtifactBucket()
			if err != nil {
				return err
			}
			bs, err := common.NewBlobStorage(bucket)
			if err != nil {
				return err
			}

			builds, err := func() ([]lib.Build, error) {
				if servicePaths == nil {
					return lib.BuildOutdatedServices(bs)
				}
				return lib.BuildServices(bs, servicePaths)
			}()
			if err != nil {
				return err
			}

			headers := []string{"service", "artifact"}
			var data [][]string
			for _, row := range builds {
				data = append(data, []string{row.Service.Path, row.ArtifactPath})
			}
			display.Table(headers, data)

			return nil
		},
	}

	cmd.Flags().StringArrayVar(&servicePaths, servicesFlag, nil, servicesDescription)
	cmd.Flags().BoolVar(&clean, cleanFlag, false, cleanDescription)

	return cmd
}
