package commands

import (
	"github.com/spf13/cobra"

	"github.com/acciaioli/mono/cmd/display"
	"github.com/acciaioli/mono/cmd/env"
	"github.com/acciaioli/mono/internal/common"
	"github.com/acciaioli/mono/lib"
	"github.com/acciaioli/mono/lib/checksum"
)

func Checksum() *cobra.Command {
	const (
		commandUse         = "checksum"
		commandDescription = "Computes/Fetches a servicePath checksum"

		serviceFlag        = "service"
		serviceDescription = "relative path to the target servicePath root directory"

		pushedFlag        = "pushed"
		pushedDescription = "fetch checksum of the latest pushed servicePath artifact"
	)

	var servicePath string
	var pushed bool

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
			service, err := lib.LoadService(servicePath)
			if err != nil {
				return err
			}
			chsum, err := func() (*string, error) {
				if pushed {
					return checksum.GetLatestChecksum(service, bs)
				}
				return checksum.ComputeChecksum(service)
			}()
			if err != nil {
				return err
			}

			display.String(*chsum)
			return nil
		},
	}

	cmd.Flags().StringVar(&servicePath, serviceFlag, "", serviceDescription)
	cmd.Flags().BoolVar(&pushed, pushedFlag, false, pushedDescription)

	return cmd
}
