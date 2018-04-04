package main

import (
	"os"

	"github.com/lordofthejars/testhubclient/hub"
	"github.com/spf13/cobra"
)

var options hub.Options

var RootCmd = &cobra.Command{
	Use:   "testhubclient",
	Short: "Interact with Test Hub",
}

func main() {

	var cmdPush = &cobra.Command{
		Use:   "push",
		Short: "Push Test Report Artifacts to Test Hub server",
		Long:  `push is used to upload test report artifacts to Test Hub server`,
		Run: func(cmd *cobra.Command, args []string) {
			hub.PublishTestReport(options, "target/surefire-reports")
		},
	}

	RootCmd.PersistentFlags().StringVarP(&options.URL, "url", "u", "http://localhost:8000", "URL where Test Hub server is deployed")
	RootCmd.PersistentFlags().StringVarP(&options.Project, "project", "p", "", "Sets Project name")
	RootCmd.PersistentFlags().StringVarP(&options.Build, "build", "b", "", "Sets Build identifier")
	RootCmd.MarkFlagRequired("project")
	RootCmd.MarkFlagRequired("build")

	RootCmd.AddCommand(cmdPush)

	if err := RootCmd.Execute(); err != nil {
		hub.Error(err.Error())
		os.Exit(1)
	}
}
