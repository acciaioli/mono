package main

import (
	"fmt"
	"os"

	"github.com/acciaioli/mono/cmd/commands"

	"github.com/spf13/cobra"
)

func main() {
	if err := execute(); err != nil {
		fmt.Printf("[error] something went wrong...\n==> %s\n", err.Error())
		os.Exit(1)
	}
}

var version = "snapshot" // build-time variable

func execute() error {
	const (
		use   = `mono`
		short = `monorepo management cli`
	)

	cmd := &cobra.Command{
		Use:           use,
		Short:         short,
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(commands.List())
	cmd.AddCommand(commands.Checksum())
	cmd.AddCommand(commands.Build())
	cmd.AddCommand(commands.Push())

	return cmd.Execute()
}
