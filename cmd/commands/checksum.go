package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/acciaioli/mono/cmd/env"
	"github.com/acciaioli/mono/services/checksum"
)

func Checksum() *cobra.Command {
	const (
		commandUse         = "checksum"
		commandDescription = "Computes/Fetches a service checksum"

		serviceFlag        = "service"
		serviceDescription = "relative path to the target service root directory"

		pushedFlag        = "pushed"
		pushedDescription = "fetch checksum of the latest pushed service artifact"
	)

	var service string
	var pushed bool

	cmd := &cobra.Command{
		Use:   commandUse,
		Short: commandDescription,
		RunE: func(cmd2 *cobra.Command, args []string) error {
			bucket, err := env.LoadArtifactBucket()
			if err != nil {
				return err
			}

			chsum, err := func() (*string, error){
				if pushed {
					return checksum.GetLatestChecksum(service, bucket)
				}
				return checksum.ComputeChecksum(service)
			}()
			if err != nil {
				return err
			}

			fmt.Println(*chsum)
			return nil
		},
	}

	cmd.Flags().StringVar(&service, serviceFlag, "", serviceDescription)
	cmd.Flags().BoolVar(&pushed, pushedFlag, false, pushedDescription)

	return cmd
}
