package commands

import (
	"github.com/spf13/cobra"

	"github.com/acciaioli/mono/cmd/display"
	"github.com/acciaioli/mono/lib"
)

func Checksum() *cobra.Command {
	const (
		commandUse         = "checksum"
		commandDescription = "Computes/Fetches a servicePath checksum"

		servicesFlag        = "service"
		servicesDescription = "relative path to the target servicePath root directory"
	)

	var servicePaths []string

	cmd := &cobra.Command{
		Use:   commandUse,
		Short: commandDescription,
		RunE: func(cmd2 *cobra.Command, args []string) error {
			checksums, err := lib.ComputeChecksums(servicePaths)
			if err != nil {
				return err
			}

			headers := []string{"service", "checksum"}
			var data [][]string
			for _, row := range checksums {
				data = append(data, []string{
					row.Service.Path,
					row.Checksum,
				})
			}
			display.Table(headers, data)

			return nil
		},
	}

	cmd.Flags().StringArrayVar(&servicePaths, servicesFlag, nil, servicesDescription)

	return cmd
}
