/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	//commands "github.com/edwinbustillos/s3tool/scripts"
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var completeCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate the autocompletion script for the specified shell",
}

func init() {
	completeCmd.Hidden = true
	rootCmd.AddCommand(completeCmd)
}
