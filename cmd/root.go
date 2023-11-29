/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var HeadAscii string = `
      ____ _______   V1.0    __
 ___ |_  //_  ___/__   __   / /
(_-<  _<_/ / /  / _ \/ _ \ / /__
/___//___/ \__/ \___/\___//____/`

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "s3tool",
	Short: "S3 tool for list,delete,upload and download files for Bucket S3",
	Long:  HeadAscii,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	var cfgFile string = "config.yaml"
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		fmt.Println("\nConfig file not found....")
		configData := `KEY: your-key
SECRET: your-secret
REGION: us-east-1`
		err := os.WriteFile(cfgFile, []byte(configData), 0644)
		if err != nil {
			fmt.Println("\nError create file config.yaml in folder: " + err.Error() + "\n")
			os.Exit(0)
		}
		fmt.Println("\nNew file created:", cfgFile+"\n")
	}
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error read file config.yaml: " + err.Error())
		os.Exit(0)
	}
}
