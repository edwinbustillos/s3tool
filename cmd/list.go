/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	commands "github.com/edwinbustillos/s3tool/scripts"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "'list' or 'ls' for bucket/files.",
	Long:    HeadAscii,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(HeadAscii)
		index := 0
		if index >= 0 && index < len(args) {
			commands.ListFilesInBucket(args)
		} else {
			commands.ListBuckets()
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
