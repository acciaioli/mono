package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/acciaioli/mono/cmd/env"
	"github.com/acciaioli/mono/services/build"
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

	var services []string
	var clean bool

	cmd := &cobra.Command{
		Use:   commandUse,
		Short: commandDescription,
		RunE: func(cmd2 *cobra.Command, args []string) error {
			if clean {
				return build.Clean()
			}

			var artifacts []build.Artifact
			var err error

			if services == nil {
				bucket, err := env.LoadArtifactBucket()
				if err != nil {
					return err
				}
				artifacts, err = build.BuildServicesWithDiff(bucket)
				if err != nil {
					return err
				}
			} else {
				artifacts, err = build.BuildServices(services)
				if err != nil {
					return err
				}
			}

			// todo: proper display
			for _, artifact := range artifacts {
				fmt.Printf("Service: %s\n", artifact.Service)
				fmt.Printf("Artifact: %s\n", artifact.Artifact)
				fmt.Printf("\n")
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVar(&services, servicesFlag, nil, servicesDescription)
	cmd.Flags().BoolVar(&clean, cleanFlag, false, cleanDescription)

	return cmd
}
