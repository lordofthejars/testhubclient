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

	var cmdDelete = &cobra.Command{
		Use:   "delete",
		Short: "Delete Specified Build from Test Hub server",
		Long:  `delete is used to remove build from Test Hub server`,
		Run: func(cmd *cobra.Command, args []string) {
			hub.RemoveBuild(options)
		},
	}

	RootCmd.PersistentFlags().StringVarP(&options.URL, "url", "u", "http://localhost:8000", "URL where Test Hub server is deployed")
	RootCmd.PersistentFlags().StringVarP(&options.Project, "project", "p", "", "Sets Project name")
	RootCmd.PersistentFlags().StringVarP(&options.Build, "build", "b", "", "Sets Build identifier")
	RootCmd.MarkFlagRequired("project")
	RootCmd.MarkFlagRequired("build")

	cmdPush.PersistentFlags().StringVarP(&options.BuildURL, "build-url", "", "", "URL where the project is built. Used for navigating from test report to build system")
	cmdPush.PersistentFlags().StringVarP(&options.Commit, "commit", "c", "", "Commit hash of current build. Used for navigating from test report to commit")
	cmdPush.PersistentFlags().StringVarP(&options.Branch, "branch", "", "", "Branch of current build. Used for navigating from test report to branch")
	cmdPush.PersistentFlags().StringVarP(&options.RepoURL, "repo-url", "r", "", "SCM location of the project. Used for navigating from test report to original source code")
	cmdPush.PersistentFlags().StringVarP(&options.RepoType, "repo-type", "t", "", "Repository type is automatically from build-url parameter. But you can explicitely set using this attribute. [github, gitlab, gogs, bitbucket]")

	RootCmd.AddCommand(cmdPush)
	RootCmd.AddCommand(cmdDelete)

	if err := RootCmd.Execute(); err != nil {
		hub.Error(err.Error())
		os.Exit(1)
	}
}
