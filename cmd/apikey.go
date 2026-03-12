/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"luna/config"
)

// apikeyCmd represents the apikey command
var apikeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Register your ApiKey from Gemini.",
	Aliases: []string{"key"},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var apiKeyStr string = args[0]

		var err error = config.SaveGlobalApiKey(apiKeyStr)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving API key: %v\n", err)
			return
		}
		fmt.Println("✅ API key saved successfully in global configuration!")
		fmt.Println("📍 Location: ~/.lunarc")

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
