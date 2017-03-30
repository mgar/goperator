package cmd

import "github.com/spf13/cobra"

var cli = &cobra.Command{
	Use:  "goperator",
	Long: "goperator allows you to manage your AWS resources",
}

// Run : RootCommand
func Run() {

	var serviceName string

	cmdListInstances.Flags().StringVarP(&serviceName, "service", "s", "", "Filter by service name")

	cli.AddCommand(cmdListInstances)
	cli.AddCommand(cmdSSHInstance)
	cli.AddCommand(execCommandInInstance)
	cli.AddCommand(cmdStopInstance)
	cli.AddCommand(cmdStartInstance)
	cli.AddCommand(cmdTerminateInstance)

	cli.Execute()
}
