package commands

import (
	"strconv"

	"github.com/acciaioli/mono/lib"

	"github.com/spf13/cobra"

	"github.com/acciaioli/mono/cmd/display"
	"github.com/acciaioli/mono/cmd/env"
	"github.com/acciaioli/mono/internal/common"
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

			listed, err := lib.List(bs)
			if err != nil {
				return err
			}

			headers := []string{"service", "diff", "version", "checksum", "local checksum"}
			var data [][]string
			for _, row := range listed {
				data = append(data, []string{
					row.Service.Path,
					strconv.FormatBool(row.Diff()),
					strconv.Itoa(row.Version),
					row.Checksum,
					row.LocalChecksum,
				})
			}
			display.Table(headers, data)

			return nil
		},
	}
	return cmd
}
