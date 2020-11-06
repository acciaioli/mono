package commands

import (
	"fmt"
	"github.com/acciaioli/mono/cmd/env"
	"github.com/spf13/cobra"

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

			// todo: proper display
			for _, service := range services {
				fmt.Printf("Service: %s\n", service.Path)
				fmt.Printf("Status: %s\n", toString(service.Status))
				fmt.Printf("Local Checksum: %s\n", toString(service.Checksum))
				fmt.Printf("Pushed Checksum: %s\n", toString(service.LatestPushedChecksum))
				fmt.Printf("\n")
			}
			return nil
		},
	}
	return cmd
}


func toString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
