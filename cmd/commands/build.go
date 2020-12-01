package commands

import (
	"github.com/spf13/cobra"

	"github.com/acciaioli/mono/cmd/display"
	"github.com/acciaioli/mono/cmd/env"
	"github.com/acciaioli/mono/internal/common"
	"github.com/acciaioli/mono/lib"
	"github.com/acciaioli/mono/lib/build"
)

func Build() *cobra.Command {
	const (
		commandUse         = "build"
		commandDescription = "Builds artifact for a service"

		servicesFlag        = "service"
		servicesDescription = "relative path to the target service root directory"

		cleanFlag        = "clean"
		cleanDescription = "cleanup builds directory"
	)

	var servicePaths []string
	var clean bool

	cmd := &cobra.Command{
		Use:   commandUse,
		Short: commandDescription,
		RunE: func(cmd2 *cobra.Command, args []string) error {
			if clean {
				return build.Clean()
			}

			var bServices []build.Service
			var err error
			bucket, err := env.LoadArtifactBucket()
			if err != nil {
				return err
			}
			bs, err := common.NewBlobStorage(bucket)
			if err != nil {
				return err
			}

			if servicePaths == nil {
				bServices, err = build.BuildOutdatedServices(bs)
				if err != nil {
					return err
				}
			} else {
				services, err := lib.LoadServices(servicePaths)
				if err != nil {
					return err
				}
				bServices, err = build.BuildServices(services, bs)
				if err != nil {
					return err
				}
			}

			headers := []string{"service", "artifact"}
			var data [][]string
			for _, bService := range bServices {
				data = append(data, []string{bService.Service.Path, bService.ArtifactPath})
			}
			display.Table(headers, data)

			return nil
		},
	}

	cmd.Flags().StringArrayVar(&servicePaths, servicesFlag, nil, servicesDescription)
	cmd.Flags().BoolVar(&clean, cleanFlag, false, cleanDescription)

	return cmd
}
