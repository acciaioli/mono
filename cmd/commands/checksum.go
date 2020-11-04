package commands

import (
	"fmt"

	"github.com/acciaioli/mono/services/checksum"
	"github.com/spf13/cobra"
)

func Checksum() *cobra.Command {
	const (
		commandUse         = "checksum"
		commandDescription = "Computes/Fetches a service checksum"

		serviceFlag        = "service"
		serviceDescription = "relative path to the target service root directory"

		pushedFlag        = "pushed"
		pushedDescription = "fetch checksum of the latest pushed service artifact"

		bucketFlag        = "bucket"
		bucketDescription = "artifacts bucket. format should be one of ['s3://',]"
	)

	var service string
	var pushed bool
	var bucket string

	cmd := &cobra.Command{
		Use:   commandUse,
		Short: commandDescription,
		RunE: func(cmd2 *cobra.Command, args []string) error {
			var chsum *string
			var err error
			if pushed {
				chsum, err = checksum.GetLatestChecksum(service, bucket)
			} else {
				chsum, err = checksum.ComputeChecksum(service)
			}
			if err != nil {
				return err
			}

			fmt.Println(*chsum)
			return nil
		},
	}

	cmd.Flags().StringVar(&service, serviceFlag, "", serviceDescription)
	cmd.Flags().BoolVar(&pushed, pushedFlag, false, pushedDescription)
	cmd.Flags().StringVar(&bucket, bucketFlag, "", bucketDescription)

	return cmd
}
