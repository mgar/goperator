package cmd

import "github.com/spf13/cobra"

var cli = &cobra.Command{
	Use:  "goperator",
	Long: "goperator allows you to manage your AWS resources",
}

// Run : RootCommand
func Run() {
	cli.AddCommand(cmdListInstances)
	cli.AddCommand(cmdSSHInstance)

	cli.Execute()
}