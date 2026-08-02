/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"luna/config"
	"luna/style"
)

// apikeyCmd represents the apikey command
var apikeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Register your OpenRouter API key.",
	Aliases: []string{"key"},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var apiKeyStr string = args[0]

		var err error = config.SaveGlobalApiKey(apiKeyStr)

		if err != nil {
			fmt.Fprintln(os.Stderr, style.Err(fmt.Sprintf("Error saving API key: %v", err)))
			return
		}
		fmt.Println(style.Success("API key saved to global configuration"))
		fmt.Println(style.Muted(style.IconMarker + "  ~/.lunarc"))

	},
}

func init() {
	rootCmd.AddCommand(apikeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apikeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apikeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
