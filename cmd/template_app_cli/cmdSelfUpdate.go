package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.melnyk.org/selfupdate-test/internal/selfupdate"
)

var selfupdateCmd = &cobra.Command{
	Use:   "selfupdate",
	Short: "SelfUpdate operation",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var selfupdateCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check available update",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		latest, err := selfupdate.GetLatestVersion(giturl)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println("Current version: ", version)
		fmt.Println("Available version: ", latest)
	},
}

var selfupdateDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download update",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	selfupdateCmd.AddCommand(selfupdateCheckCmd)
	selfupdateCmd.AddCommand(selfupdateDownloadCmd)
	rootCmd.AddCommand(selfupdateCmd)
}
