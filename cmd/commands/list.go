package commands

import (
	"github.com/spf13/cobra"

	"github.com/acciaioli/mono/cmd/display"
	"github.com/acciaioli/mono/cmd/env"
	"github.com/acciaioli/mono/services/list"
)

func List() *cobra.Command {
	const (
		commandUse         = "list"
		commandDescription = "Lists all services under the current directory"
	)

	cmd := &cobra.Command{
		Use:   commandUse,
		Short: commandDescription,
		RunE: func(cmd2 *cobra.Command, args []string) error {
			bucket, err := env.LoadArtifactBucket()
			if err != nil {
				return err
			}

			services, err := list.List(bucket)
			if err != nil {
				return err
			}

			headers := []string{"service", "status", "checksum", "local checksum"}
			var data [][]string
			for _, service := range services {
				data = append(data, []string{service.Path, string(service.Status), service.LatestPushedChecksum, service.Checksum})
			}
			display.Table(headers, data)

			return nil
		},
	}
	return cmd
}
