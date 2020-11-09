package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/acciaioli/mono/services/build"
)

func Build() *cobra.Command {
	const (
		commandUse         = "build"
		commandDescription = "Builds artifact for a service"

		serviceFlag        = "service"
		serviceDescription = "relative path to the target service root directory"
	)

	var service string

	cmd := &cobra.Command{
		Use:   commandUse,
		Short: commandDescription,
		RunE: func(cmd2 *cobra.Command, args []string) error {
			artifact, err := build.Build(service)
			if err != nil {
				return err
			}

			fmt.Println(*artifact)

			return nil
		},
	}

	cmd.Flags().StringVar(&service, serviceFlag, "", serviceDescription)

	return cmd
}
