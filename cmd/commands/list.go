package commands

import (
	"github.com/spf13/cobra"

	"github.com/acciaioli/mono/cmd/display"
	"github.com/acciaioli/mono/cmd/env"
	"github.com/acciaioli/mono/internal/common"
	"github.com/acciaioli/mono/lib/list"
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
			bs, err := common.NewBlobStorage(bucket)
			if err != nil {
				return err
			}

			lServices, err := list.List(bs)
			if err != nil {
				return err
			}

			headers := []string{"service", "status", "checksum", "local checksum"}
			var data [][]string
			for _, lService := range lServices {
				data = append(data, []string{
					lService.Service.Path,
					string(lService.Status),
					lService.LatestPushedChecksum,
					lService.Checksum,
				})
			}
			display.Table(headers, data)

			return nil
		},
	}
	return cmd
}
