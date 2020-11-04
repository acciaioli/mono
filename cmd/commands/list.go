package commands

import (
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
			return list.List()
		},
	}

	return cmd
}
