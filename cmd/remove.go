/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	commands "github.com/edwinbustillos/s3tool/scripts"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "'remove' or 'rm' for delete files in Bucket",
	Long:    `'remove' or 'rm' for delete files in Bucket E.g: bucket-name/path/file.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(HeadAscii)
		index := 0
		if index >= 0 && index < len(args) {
			commands.RemoveFolder(args)
		} else {
			fmt.Println("\nError path file, try again E.g: s3tool rm bucketName/path/file.txt")
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolP("folder", "f", true, "Flag for remove folder and files recursive.")
}
